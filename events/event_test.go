package events

import (
	"strings"
	"testing"
	"time"
)

func TestParseEvent(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		wantEvent Event
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "Empty line",
			input:   "",
			wantErr: true,
			errMsg:  "пустая строка",
		},
		{
			name:    "Registration event",
			input:   "[09:00:00.000] 1 1",
			wantErr: false,
			wantEvent: Event{
				Time:      time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
				EventID:   EventRegister,
				AthleteID: 1,
				Params:    []string{},
				Raw:       "[09:00:00.000] 1 1",
			},
		},
		{
			name:    "Start time lottery",
			input:   "[09:05:00.000] 2 1 09:30:00.000",
			wantErr: false,
			wantEvent: Event{
				Time:      time.Date(0, 1, 1, 9, 5, 0, 0, time.UTC),
				EventID:   EventStartTimeLottery,
				AthleteID: 1,
				Params:    []string{"09:30:00.000"},
				Raw:       "[09:05:00.000] 2 1 09:30:00.000",
			},
		},
		{
			name:    "At firing line with params",
			input:   "[09:45:00.000] 5 1 1",
			wantErr: false,
			wantEvent: Event{
				Time:      time.Date(0, 1, 1, 9, 45, 0, 0, time.UTC),
				EventID:   EventAtFiringLine,
				AthleteID: 1,
				Params:    []string{"1"},
				Raw:       "[09:45:00.000] 5 1 1",
			},
		},
		{
			name:    "Hit successful",
			input:   "[09:45:05.000] 6 1 1",
			wantErr: false,
			wantEvent: Event{
				Time:      time.Date(0, 1, 1, 9, 45, 5, 0, time.UTC),
				EventID:   EventHitSuccessful,
				AthleteID: 1,
				Params:    []string{"1"},
				Raw:       "[09:45:05.000] 6 1 1",
			},
		},
		{
			name:    "Hit missed",
			input:   "[09:45:06.000] 61 1 2",
			wantErr: false,
			wantEvent: Event{
				Time:      time.Date(0, 1, 1, 9, 45, 6, 0, time.UTC),
				EventID:   EventHitMissed,
				AthleteID: 1,
				Params:    []string{"2"},
				Raw:       "[09:45:06.000] 61 1 2",
			},
		},
		{
			name:    "Invalid event format",
			input:   "[09:45:06.000] invalid",
			wantErr: true,
			errMsg:  "Некорректная строка события: [09:45:06.000] invalid",
		},
		{
			name:    "Invalid time format",
			input:   "[09:45:06] 1 1",
			wantErr: true,
			errMsg:  "Ошибка парсинга времени",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseEvent(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseEvent() expected error, got nil")
				} else if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("ParseEvent() error = %v, want error containing %q", err, tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseEvent() unexpected error = %v", err)
				return
			}

			if !got.Time.Equal(tc.wantEvent.Time) {
				t.Errorf("ParseEvent() Time = %v, want %v", got.Time, tc.wantEvent.Time)
			}

			if got.EventID != tc.wantEvent.EventID {
				t.Errorf("ParseEvent() EventID = %d, want %d", got.EventID, tc.wantEvent.EventID)
			}

			if got.AthleteID != tc.wantEvent.AthleteID {
				t.Errorf("ParseEvent() AthleteID = %d, want %d", got.AthleteID, tc.wantEvent.AthleteID)
			}

			if !equalStringSlices(got.Params, tc.wantEvent.Params) {
				t.Errorf("ParseEvent() Params = %v, want %v", got.Params, tc.wantEvent.Params)
			}

			if got.Raw != tc.wantEvent.Raw {
				t.Errorf("ParseEvent() Raw = %q, want %q", got.Raw, tc.wantEvent.Raw)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
