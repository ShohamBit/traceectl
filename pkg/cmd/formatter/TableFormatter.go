package formatter

import (
	"fmt"
	"os"
	"strings"

	"github.com/aquasecurity/table"
	pb "github.com/aquasecurity/tracee/api/v1beta1"
)

func (f *Formatter) PrintSteamTableHeaders() {
	f.CMD.Printf("%-15s %-25s %-15s %-15s %s\n",
		"TIME",
		"EVENT NAME",
		"POLICIES",
		"PID",
		"DATA",
	)
}
func (f *Formatter) PrintStreamTableRow(event *pb.Event) {
	timestamp := event.Timestamp.AsTime().Format("15:04:05.000")

	f.CMD.Printf("%-15s %-25s %-15s %-15s %s\n",
		timestamp,
		event.Name,
		strings.Join(event.Policies.Matched, ","),
		fmt.Sprintf("%d", event.Context.Process.Pid.Value),
		getEventData(event.Data),
	)

}

// generate event data
func getEventData(data []*pb.EventValue) string {
	var result []string
	for _, ev := range data {
		result = append(result, getEventName(ev)+getEventValue(ev))
	}
	return strings.Join(result, ", ")
}
func getEventName(ev *pb.EventValue) string {
	return strings.ToUpper(ev.Name[0:1]) + ev.Name[1:] + ": "
}

// generate event value
func getEventValue(ev *pb.EventValue) string {
	switch v := ev.Value.(type) {
	case *pb.EventValue_Int32:
		return fmt.Sprintf("%d", v.Int32)
	case *pb.EventValue_Int64:
		return fmt.Sprintf("%d", v.Int64)
	case *pb.EventValue_UInt32:
		return fmt.Sprintf("%d", v.UInt32)
	case *pb.EventValue_UInt64:
		return fmt.Sprintf("%d", v.UInt64)
	case *pb.EventValue_Str:
		return v.Str
	case *pb.EventValue_Bytes:
		return fmt.Sprintf("%x", v.Bytes)
	case *pb.EventValue_Bool:
		return fmt.Sprintf("%t", v.Bool)
	case *pb.EventValue_StrArray:
		return strings.Join(v.StrArray.Value, ", ")
	case *pb.EventValue_Int32Array:
		return fmt.Sprintf("%v", v.Int32Array.Value)
	case *pb.EventValue_UInt64Array:
		return fmt.Sprintf("%v", v.UInt64Array.Value)
		//TODO: add more types
	default:
		// if data type not supported yet
		return "unknown"
	}
}

func (f *Formatter) PrintEventListTable(response *pb.GetEventDefinitionsResponse) *table.Table {
	tbl := table.New(os.Stdout)
	tbl.SetHeaders("ID", "Name", "Version", "Tags")
	for _, event := range response.Definitions {
		// Check if the optional field Threat is set (non-nil)

		tbl.AddRow(
			fmt.Sprintf("%d", event.Id),
			event.Name,
			fmt.Sprintf("%d.%d.%d", event.Version.Major, event.Version.Minor, event.Version.Patch),
			strings.Join(event.Tags, ", "),
		)

	}
	return tbl
}

func (f *Formatter) PrintEventDescriptionTable(response *pb.GetEventDefinitionsResponse) *table.Table {
	tbl := table.New(os.Stdout)
	tbl.SetHeaders("ID", "Name", "Version", "Tags", "Description")
	for _, event := range response.Definitions {
		// Check if the optional field Threat is set (non-nil)

		tbl.AddRow(
			fmt.Sprintf("%d", event.Id),
			event.Name,
			fmt.Sprintf("%d.%d.%d", event.Version.Major, event.Version.Minor, event.Version.Patch),
			strings.Join(event.Tags, ", "),
			event.Description,
		)

	}
	return tbl
}
