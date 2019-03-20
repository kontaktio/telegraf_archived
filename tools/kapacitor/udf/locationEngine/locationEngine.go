package locationEngine

import (
	"log"
	"os"

	"github.com/influxdata/kapacitor/udf/agent"
)

const cRssiThreshold = 3

// Location handler template
type locationHandler struct {
	window  *locationWindow
	begin   *agent.BeginBatch
	agent   *agent.Agent
	history *locationWindow
}

//map of tracking id's
type locationWindow struct {
	entries map[string]*locationReference
}

type locationReference struct {
	rssi     float64
	sourceID string
	point    *agent.Point
}

func (o *locationReference) Clone() *locationReference {
	obj := &locationReference{rssi: o.rssi, sourceID: o.sourceID, point: o.point}
	return obj
}

func newlocationHandler(agent *agent.Agent) *locationHandler {
	historyWindow := &locationWindow{}
	historyWindow.reset()
	return &locationHandler{agent: agent, window: &locationWindow{}, history: historyWindow}
}

// Return the InfoResponse. Describing the properties of this UDF agent.
func (*locationHandler) Info() (*agent.InfoResponse, error) {
	info := &agent.InfoResponse{
		Wants:    agent.EdgeType_BATCH,
		Provides: agent.EdgeType_BATCH,
		Options:  map[string]*agent.OptionInfo{},
	}
	return info, nil
}

// Initialze the handler based of the provided options.
func (o *locationHandler) Init(r *agent.InitRequest) (*agent.InitResponse, error) {
	init := &agent.InitResponse{
		Success: true,
		Error:   "",
	}
	return init, nil
}

// Create a snapshot of the running window of the process.
func (o *locationHandler) Snapshot() (*agent.SnapshotResponse, error) {

	return &agent.SnapshotResponse{
		Snapshot: nil,
	}, nil
}

// Restore a previous snapshot.
func (o *locationHandler) Restore(req *agent.RestoreRequest) (*agent.RestoreResponse, error) {

	return &agent.RestoreResponse{
		Success: true,
		Error:   "",
	}, nil
}

// Start working with the next batch
func (o *locationHandler) BeginBatch(begin *agent.BeginBatch) error {
	o.window.reset()
	// Keep begin batch for later
	o.begin = begin
	return nil
}

// Add point to map
func (o *locationHandler) Point(p *agent.Point) error {
	trackingId := p.Tags["trackingId"]
	sourceId, exists := p.FieldsString["sourceId"]
	if !exists {
		sourceId = p.Tags["sourceId"]
	}
	location := o.window.entries[trackingId]
	newRssi, ok := p.FieldsDouble["rssi"]
	if ok {
		if location == nil {

			o.window.entries[trackingId] = &locationReference{rssi: newRssi, sourceID: sourceId, point: p}
		} else if location.rssi < newRssi {

			o.window.entries[trackingId] = &locationReference{rssi: newRssi, sourceID: sourceId, point: p}
		} else {

		}

	}

	return nil
}

func adjustZonesBasedOnHistor(o *locationHandler) {
	//create new history data
	newHistory := &locationWindow{}
	newHistory.reset()
	currentEntries := o.window.entries
	historyEntries := o.history.entries

	for k, entry := range currentEntries {
		//save current data as new historical data
		newHistory.entries[k] = currentEntries[k].Clone()

		//for each sourceId in current data compare with history data
		//If new point has different source id and is not stronger that 3 dbm
		//than use historical zone data with new timestamp. Else use current one.
		historyEntry, exists := historyEntries[k]

		if (exists) &&
			(historyEntry.sourceID != entry.sourceID) &&
			(historyEntry.rssi >= entry.rssi-cRssiThreshold) {

			currentEntries[k].rssi = historyEntry.rssi
			currentEntries[k].sourceID = historyEntry.sourceID
		} else {
			//entry not in history
			continue
		}

	}

	o.history = newHistory
}

// Finish batch and get calculated points
func (o *locationHandler) EndBatch(end *agent.EndBatch) error {

	o.begin.Size = int64(len(o.window.entries))
	// End batch
	o.agent.Responses <- &agent.Response{
		Message: &agent.Response_Begin{
			Begin: o.begin,
		},
	}

	adjustZonesBasedOnHistor(o)

	entries := o.window.getPoints()
	for _, location := range entries {
		o.agent.Responses <- &agent.Response{
			Message: &agent.Response_Point{
				Point: location,
			},
		}
	}

	log.Printf("LocationEngine batch completed with %d points", len(entries))
	o.agent.Responses <- &agent.Response{
		Message: &agent.Response_End{
			End: end,
		},
	}
	return nil
}

// Stop the handler gracefully.
func (o *locationHandler) Stop() {
	close(o.agent.Responses)
}

// create database entries from points
func (w *locationWindow) getPoints() []*agent.Point {
	entries := len(w.entries)
	points := make([]*agent.Point, 0, entries)

	for k, entry := range w.entries {
		point := entry.point
		point.FieldsInt = nil
		point.FieldsBool = nil
		point.FieldsDouble = map[string]float64{"rssi": entry.rssi}
		point.FieldsString = map[string]string{"fTrackingId": k, "fSourceId": entry.sourceID}
		point.Tags = map[string]string{"trackingId": k, "sourceId": entry.sourceID}
		points = append(points, point)
	}

	return points
}

func (w *locationWindow) reset() {
	w.entries = make(map[string]*locationReference)
}

func main() {
	a := agent.New(os.Stdin, os.Stdout)
	h := newlocationHandler(a)
	a.Handler = h

	log.Println("Starting locationEngine agent")
	a.Start()
	log.Println("locationEngine agent started")
	err := a.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
