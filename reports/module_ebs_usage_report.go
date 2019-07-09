//   Copyright 2019 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package reports

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/ebs"
	"github.com/trackit/trackit-server/users"
)

const ebsUsageReportSheetName = "EBS Usage Report"

var ebsUsageReportModule = module{
	Name:          "EBS Usage Report",
	SheetName:     ebsUsageReportSheetName,
	ErrorName:     "ebsUsageReportError",
	GenerateSheet: generateEbsUsageReportSheet,
}

func generateEbsUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return ebsUsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func ebsUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := ebsUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return ebsUsageReportInsertDataInSheet(ctx, aas, file, data)
	} else {
		return
	}
}

func ebsUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []ebs.SnapshotReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	identities := getAwsIdentities(aas)

	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}

	parameters := ebs.EbsQueryParams{
		AccountList: identities,
		Date:        date,
	}

	logger.Debug("Getting EBS Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = ebs.GetEbsData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an EBS Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func ebsUsageReportInsertDataInSheet(ctx context.Context, aas []aws.AwsAccount, file *excelize.File, data []ebs.SnapshotReport) (err error) {
	file.NewSheet(ebsUsageReportSheetName)
	ebsUsageReportGenerateHeader(file)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	line := 3
	logger.Debug("Getting informations from EBS, TRYING TO GET SOME INFORMATIONS", map[string]interface{}{
		"data": data,
	})
	for _, report := range data {
		logger.Debug("Getting informations from EBS, TRYING TO LOOP ON DATA", "")
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		snapshot := report.Snapshot
		date := snapshot.StartTime.Format("2006-01-02 15:04:05")
		tags := formatTags(snapshot.Tags)
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(snapshot.Id, "B"+strconv.Itoa(line)),
			newCell(date, "C"+strconv.Itoa(line)),
			newCell(snapshot.Region, "D"+strconv.Itoa(line)),
			newCell(snapshot.Cost, "E"+strconv.Itoa(line)),
			newCell(snapshot.State, "F"+strconv.Itoa(line)),
			newCell(snapshot.Volume.Id, "G"+strconv.Itoa(line)),
			newCell(snapshot.Volume.Size, "H"+strconv.Itoa(line)),
			newCell(strings.Join(tags, ";"), "I"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, ebsUsageReportSheetName)
		line++
	}
	return
}

func ebsUsageReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("ID", "B1").mergeTo("B2"),
		newCell("Date", "C1").mergeTo("C2"),
		newCell("Region", "D1").mergeTo("D2"),
		newCell("Cost", "E1").mergeTo("E2"),
		newCell("State", "F1").mergeTo("F2"),
		newCell("Volume", "G1").mergeTo("H1"),
		newCell("ID", "G2"),
		newCell("Size (GigaBytes)", "H2"),
		newCell("Tags", "I1").mergeTo("I2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, ebsUsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 40),
		newColumnWidth("C", 20).toColumn("F"),
		newColumnWidth("G", 30),
		newColumnWidth("H", 20),
		newColumnWidth("I", 30),
	}
	columns.setValues(file, ebsUsageReportSheetName)
	return
}