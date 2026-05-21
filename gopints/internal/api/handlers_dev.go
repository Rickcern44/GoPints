package api

import (
	"net/http"
	"time"
)

func (s *Server) handleSimulatePour(w http.ResponseWriter, r *http.Request) {
	if !s.simulate {
		writeError(w, http.StatusForbidden, "not in simulate mode")
		return
	}
	tapID, err := parseTapID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tap id")
		return
	}
	pulses := queryInt(r, "pulses", 200)

	meter, ok := s.meters[tapID]
	if !ok {
		writeError(w, http.StatusNotFound, "tap not found")
		return
	}

	go func() {
		for i := 0; i < pulses; i++ {
			meter.HandlePulse()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	writeJSON(w, http.StatusAccepted, map[string]any{
		"tap_id": tapID,
		"pulses": pulses,
	})
}
