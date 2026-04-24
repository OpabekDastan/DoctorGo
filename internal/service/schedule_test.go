package service

import (
	"testing"
)

func TestValidateScheduleInput(t *testing.T) {
	svc := &ScheduleService{}

	tests := []struct {
		name    string
		input   UpsertScheduleInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: UpsertScheduleInput{
				Weekday:     1,
				StartTime:   "09:00",
				EndTime:     "17:00",
				SlotMinutes: 30,
			},
			wantErr: false,
		},
		{
			name: "weekday too high (8)",
			input: UpsertScheduleInput{
				Weekday:     8,
				StartTime:   "09:00",
				EndTime:     "17:00",
				SlotMinutes: 30,
			},
			wantErr: true,
		},
		{
			name: "weekday negative (-1)",
			input: UpsertScheduleInput{
				Weekday:     -1,
				StartTime:   "09:00",
				EndTime:     "17:00",
				SlotMinutes: 30,
			},
			wantErr: true,
		},
		{
			name: "missing start_time",
			input: UpsertScheduleInput{
				Weekday:     1,
				StartTime:   "",
				EndTime:     "17:00",
				SlotMinutes: 30,
			},
			wantErr: true,
		},
		{
			name: "missing end_time",
			input: UpsertScheduleInput{
				Weekday:     1,
				StartTime:   "09:00",
				EndTime:     "",
				SlotMinutes: 30,
			},
			wantErr: true,
		},
		{
			name: "zero slot_minutes",
			input: UpsertScheduleInput{
				Weekday:     1,
				StartTime:   "09:00",
				EndTime:     "17:00",
				SlotMinutes: 0,
			},
			wantErr: true,
		},
		{
			name: "negative slot_minutes",
			input: UpsertScheduleInput{
				Weekday:     1,
				StartTime:   "09:00",
				EndTime:     "17:00",
				SlotMinutes: -15,
			},
			wantErr: true,
		},
		{
			name: "boundary weekday 0 (Monday)",
			input: UpsertScheduleInput{
				Weekday:     0,
				StartTime:   "08:00",
				EndTime:     "16:00",
				SlotMinutes: 60,
			},
			wantErr: false,
		},
		{
			name: "boundary weekday 6 (Sunday)",
			input: UpsertScheduleInput{
				Weekday:     6,
				StartTime:   "10:00",
				EndTime:     "14:00",
				SlotMinutes: 30,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
