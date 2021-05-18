package main

import (
	"errors"
	"fitbot/internal/fitforfree"
	"fitbot/internal/fitforfree/lesson"
	"fitbot/internal/fitforfree/login"
	"fitbot/internal/fitforfree/scheduled_classes"
	"fitbot/internal/fitforfree/utils"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	env "github.com/caarlos0/env"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	dotenv "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	funk "github.com/thoas/go-funk"
)

var (
	botAPI *telegram.BotAPI
	config struct {
		Debug       bool   `env:"DEBUG"`
		BotOwner    int    `env:"BOT_OWNER,required"`
		BotToken    string `env:"BOT_TOKEN,required"`
		FitMemberID string `env:"FIT_MEMBER_ID,required"`
		FitPostcode string `env:"FIT_POSTCODE,required"`
		FitVenue    string `env:"FIT_VENUE,required"`
	}
	lessionAPI          lesson.API
	scheduledClassesAPI scheduled_classes.API
)

func reply(message *telegram.Message, text string) error {
	msg := telegram.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID
	if _, err := botAPI.Send(msg); err != nil {
		return fmt.Errorf("failed to send: %+v", err)
	}
	return nil
}

var (
	bookSig chan struct{}
	bookWg  sync.WaitGroup
)

func book(message *telegram.Message, params []string) error {
	if len(params) != 4 {
		return errors.New("invalid format, expected: book days start-hour end-hour iterations")
	}

	dayStrings := strings.Split(params[0], ",")
	days := make([]int, len(dayStrings))
	for i, s := range dayStrings {
		day, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("failed to parse day: %+v", err)
		}
		days[i] = day
	}
	startHour, err := strconv.Atoi(params[1])
	if err != nil {
		return fmt.Errorf("failed to parse start-hour: %+v", err)
	}
	endHour, err := strconv.Atoi(params[2])
	if err != nil {
		return fmt.Errorf("failed to parse end-hour: %+v", err)
	}
	iterations, err := strconv.Atoi(params[3])
	if err != nil {
		return fmt.Errorf("failed to parse iterations: %+v", err)
	}

	runDays := make([]int, len(days)*2-1) // -1 because the last element of days only occurs once, in the middle
	copy(runDays, days)
	copy(runDays[len(days)-1:], funk.ReverseInt(days))
	runCompleted := make([]bool, len(runDays))
	run := func() (bool, error) {
	nextDay:
		for i, day := range runDays {
			log := log.WithField("day", day)
			if runCompleted[i] {
				log.Debug("completed")
				continue nextDay
			}

			time.Sleep(time.Duration(funk.RandomInt(3*1000, 4*1000)) * time.Millisecond)

			log.Debug("listing")
			res, err := scheduledClassesAPI.List(scheduled_classes.ListQuery{
				Category:           "free_practise",
				Date:               utils.Date(day),
				ExcludeFullyBooked: true,
				Venues:             config.FitVenue,
			})
			if err != nil {
				return false, fmt.Errorf("failed to list scheduled classes: %+v", err)
			}

		nextLesson:
			for _, cur := range res.ScheduledClasses {
				if cur.ClassType != "free_practise" || cur.Activity.ID != "vrijtrainen" {
					continue nextLesson
				}

				startsAt := cur.Start().Hour()
				if startsAt < startHour || startsAt > endHour {
					continue nextLesson
				}

				if cur.Booked {
					log.WithField("lesson", cur.String()).Debug("already booked, mark completed")
					runCompleted[i] = true
					continue nextDay
				}

				if cur.Status == "AVAILABLE" {
					time.Sleep(time.Duration(funk.RandomInt(2*1000, 3*1000)) * time.Millisecond)

					log := log.WithField("lesson", cur.String())
					log.Debug("booking")
					booking, err := lessionAPI.Book(lesson.BookData{ID: cur.ID})
					if err != nil {
						return false, fmt.Errorf("failed to book: %+v", err)
					}
					log.WithField("booking", booking).Debug("booked")
					if err := reply(message, fmt.Sprintf("booked: %v", cur.String())); err != nil {
						return false, fmt.Errorf("failed to send: %+v", err)
					}
					runCompleted[i] = true
					continue nextDay
				}
			}
		}
		return !funk.Contains(runCompleted, false), nil
	}

	if bookSig != nil {
		close(bookSig)
	}
	bookWg.Wait()

	if err := reply(message, "started"); err != nil {
		return fmt.Errorf("failed to send: %+v", err)
	}

	go func() {
		bookWg.Add(1)
		bookSig = make(chan struct{})
		defer func() {
			if err := reply(message, "ended"); err != nil {
				log.WithError(err).Error("failed to send")
			}
			bookSig = nil
			bookWg.Done()
		}()

		i := 0
		delay := 0
		for {
			if i > 0 {
				delay = funk.RandomInt(5*60*1000, 10*60*1000)
			}

			select {
			case <-time.After(time.Duration(delay) * time.Millisecond):
				i++
				var text string

				done, err := run()
				if err != nil {
					text = fmt.Sprintf("failed to run: %+v", err)
				} else if done {
					text = "all days booked"
				} else if i == iterations {
					text = "iterations limit reached"
				}

				if text != "" {
					if err := reply(message, text); err != nil {
						log.WithError(err).Error("failed to send")
					}
					return
				}
			case <-bookSig:
				return
			}
		}
	}()

	return nil
}

func stop(message *telegram.Message, params []string) error {
	if bookSig == nil {
		return errors.New("not running")
	}

	close(bookSig)
	return nil
}

var handlers = map[string]func(message *telegram.Message, params []string) error{
	"book": book,
	"stop": stop,
}

func main() {
	dotenv.Load()
	if err := env.Parse(&config); err != nil {
		log.WithError(err).Fatal("failed to parse environment")
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	API := *fitforfree.New()
	lessionAPI = lesson.API{API: API}
	loginAPI := login.API{API: API}
	scheduledClassesAPI = scheduled_classes.API{API: API}

	var err error
	LoginResult, err := loginAPI.Login(login.LoginData{
		MemberID:      config.FitMemberID,
		Postcode:      config.FitPostcode,
		TermsAccepted: true,
	})
	if err != nil {
		log.WithError(err).Fatal("failed to login")
	}
	API.SetAuth(LoginResult.Data.SessionID)
	log.Infof("logged in as %s %s", LoginResult.Data.FirstName, LoginResult.Data.Surname)

	botAPI, err = telegram.NewBotAPI(config.BotToken)
	if err != nil {
		log.WithError(err).Fatal("failed to create bot")
	}

	u := telegram.NewUpdate(0)
	u.Timeout = 60
	updates, err := botAPI.GetUpdatesChan(u)
	if err != nil {
		log.WithError(err).Fatal("failed to get channel updates")
	}

	for update := range updates {
		if update.Message == nil || update.Message.From.ID != config.BotOwner {
			continue
		}

		var err error
		parts := strings.Split(update.Message.Text, " ")
		command, params := parts[0], parts[1:]
		if handler, ok := handlers[command]; ok {
			err = handler(update.Message, params)
		} else {
			err = errors.New("unknown command")
		}

		if err != nil {
			if err := reply(update.Message, fmt.Sprintf("%+v", err)); err != nil {
				log.WithError(err).Error("failed to send")
			}
		}
	}
}
