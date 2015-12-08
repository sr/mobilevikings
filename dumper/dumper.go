package dumper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/sr/mobilevikings"
)

type dumper struct {
	client    mobilevikings.Client
	directory string
}

func NewDumper(
	client mobilevikings.Client,
	directory string,
) *dumper {
	return &dumper{client, directory}
}

func (d *dumper) Dump() error {
	if fileInfo, err := os.Stat(d.directory); os.IsNotExist(err) {
		return fmt.Errorf("directory \"%s\" not found", d.directory)
	} else if !fileInfo.IsDir() {
		return fmt.Errorf("file \"%s\" is not a directory", d.directory)
	}

	phoneNumbers, err := d.client.PhoneNumbers()
	if err != nil {
		return err
	}

	for _, phoneNumber := range phoneNumbers {
		if err := d.dumpPhoneNumber(phoneNumber); err != nil {
			return err
		}
	}
	return nil
}

func (d *dumper) dumpPhoneNumber(phoneNumber mobilevikings.PhoneNumber) error {
	directory := path.Join(d.directory, phoneNumber.ID)
	usageDir := path.Join(directory, "usage")
	topupDir := path.Join(directory, "topup")
	if err := os.MkdirAll(usageDir, 0775); err != nil {
		return err
	}
	if err := d.dumpUsage(phoneNumber.ID, usageDir); err != nil {
		return err
	}
	if err := os.MkdirAll(topupDir, 0775); err != nil {
		return err
	}
	if err := d.dumpTopup(phoneNumber.ID, topupDir); err != nil {
		return err
	}
	return nil
}

func (d *dumper) dumpUsage(phoneNumberID string, directory string) error {
	insights, err := d.client.Insights(phoneNumberID)
	if err != nil {
		return err
	}
	daysAsViking := insights.VikingLife.DaysAsViking
	signup := time.Now().AddDate(0, 0, -daysAsViking)
	signupBeginningOfMonth := time.Date(
		signup.Year(),
		signup.Month(),
		1,
		00,
		00,
		00,
		0,
		time.UTC,
	)

	i := 0
	for {
		from := signupBeginningOfMonth.AddDate(0, 0+i, 0)
		until := from.AddDate(0, 1, 0).Add(-time.Minute)
		fileName := fmt.Sprintf("%d-%02d.%s", from.Year(), from.Month(), "json")
		filePath := path.Join(directory, fileName)

		usages, err := d.client.Usage(phoneNumberID, from, until)
		if err != nil {
			return err
		}
		if len(usages) == 0 {
			fmt.Printf("no usages returned")
			break
		}

		marshalled, err := json.Marshal(usages)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filePath, marshalled, 0644); err != nil {
			return err
		}

		i = i + 1
	}
	return nil
}

func (d *dumper) dumpTopup(phoneNumberID string, directory string) error {
	var topups []mobilevikings.Topup
	pageURL := ""
	for {
		page, err := d.client.Topups(phoneNumberID, pageURL)
		if err != nil {
			return err
		}
		for _, topup := range page.Results {
			topups = append(topups, topup)
		}
		if page.Next == "" {
			break
		}
		pageURL = page.Next
	}

	marshalled, err := json.Marshal(topups)
	if err != nil {
		return err
	}

	filePath := path.Join(directory, "all.json")
	return ioutil.WriteFile(filePath, marshalled, 0644)
}
