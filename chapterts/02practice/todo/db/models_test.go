package db

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

func TestBuilder_Build(t *testing.T) {
	mockId := uuid.MustParse("01964483-01b5-779f-9c6f-b2496503591d")
	testIdGenerator := func() uuid.UUID {
		return mockId
	}

	now := time.Now()

	var tests = []struct {
		testName    string
		name        string
		description string
		time        time.Time
		expected    Task
	}{
		{
			testName:    "basic task",
			name:        "Buy groceries",
			description: "Milk, eggs, bread",
			time:        now,
			expected: Task{
				ID:          mockId,
				Name:        "Buy groceries",
				Description: "Milk, eggs, bread",
				Time:        now,
			},
		},
		{
			testName:    "task with empty description",
			name:        "Walk the dog",
			description: "",
			time:        now.Add(time.Hour),
			expected: Task{
				ID:          mockId,
				Name:        "Walk the dog",
				Description: "",
				Time:        now.Add(time.Hour),
			},
		},
		{
			testName:    "task with zero time",
			name:        "Pay bills",
			description: "Electricity, water",
			time:        time.Time{},
			expected: Task{
				ID:          mockId,
				Name:        "Pay bills",
				Description: "Electricity, water",
				Time:        time.Time{},
			},
		},
		{
			testName:    "task with special characters",
			name:        "Learn Go!",
			description: "!@#$%^&*()_+",
			time:        now.AddDate(0, 1, 0),
			expected: Task{
				ID:          mockId,
				Name:        "Learn Go!",
				Description: "!@#$%^&*()_+",
				Time:        now.AddDate(0, 1, 0),
			},
		},
		{
			testName:    "task with long description",
			name:        "Write report",
			description: "This is a very long description for the report. It needs to cover all the key findings and recommendations from the last quarter's analysis. We should also include some projections for the next quarter based on the current trends. Make sure to cite all the sources properly and include a detailed appendix with all the relevant data.",
			time:        now.AddDate(0, 0, 7),
			expected: Task{
				ID:          mockId,
				Name:        "Write report",
				Description: "This is a very long description for the report. It needs to cover all the key findings and recommendations from the last quarter's analysis. We should also include some projections for the next quarter based on the current trends. Make sure to cite all the sources properly and include a detailed appendix with all the relevant data.",
				Time:        now.AddDate(0, 0, 7),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result := NewTaskBuilder(testIdGenerator).
				WithName(tt.name).
				WithDescription(tt.description).
				WithTime(tt.time).
				Build()
			if !reflect.DeepEqual(tt.expected, *result) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func setupMockFS() (afero.Fs, func()) {
	fs := afero.NewMemMapFs()
	appFs = fs
	return fs, func() {
		appFs = afero.NewOsFs()
	}
}

var testTask0 = NewTaskBuilder(UuidIdGenerator).Build()
var testTask1 = NewTaskBuilder(UuidIdGenerator).Build()

var testTasksData = map[string]*Task{
	testTask0.ID.String(): testTask0,
	testTask1.ID.String(): testTask1,
}

func TestGetDataFromFs_Success(t *testing.T) {
	fs, teardown := setupMockFS()
	defer teardown()

	jsonData, _ := json.Marshal(testTasksData)
	afero.WriteFile(fs, storageFp, jsonData, 0644)

	result, err := getDataFromFs()
	if err != nil {
		t.Errorf("GetDataFromFs returned an Error: %v", err)
	}

	if !reflect.DeepEqual(testTasksData, result) {
		t.Errorf("GetDataFromFs returned incorrect data. Got: %v, Want: %v", result, testTasksData)
	}
}

func TestGetDataFromFs_FileNotExist(t *testing.T) {
	_, teardown := setupMockFS()
	defer teardown()

	data, err := getDataFromFs()

	if err != nil && !os.IsNotExist(err) {
		t.Errorf("GetDataFromFs returned an unexpected error for a non-existent file: %v", err)
	}
	if data != nil {
		t.Errorf("GetDataFromFs should return nil data for a non-existent file. Got: %v", data)
	}
}

func TestGetDataFromFs_InvalidJSON(t *testing.T) {
	fs, teardown := setupMockFS()
	defer teardown()

	afero.WriteFile(fs, storageFp, []byte("this is not valid json"), 0644)

	data, err := getDataFromFs()

	if err == nil {
		t.Errorf("GetDataFromFs should have returned an error for invalid JSON")
	}
	if data != nil {
		t.Errorf("GetDataFromFs should return nil data for invalid JSON. Got: %v", data)
	}
}

func TestSaveDataToFs_Success(t *testing.T) {
	fs, teardown := setupMockFS()
	defer teardown()

	err := saveDataToFs(testTasksData)
	if err != nil {
		t.Errorf("saveDataToFs returned an error: %v", err)
	}

	readDataBytes, err := afero.ReadFile(fs, storageFp)
	if err != nil {
		t.Fatalf("Failed to read data from mock file system: %v", err)
	}

	var readData map[string]*Task
	err = json.Unmarshal(readDataBytes, &readData)
	if err != nil {
		t.Fatalf("Failed to unmarshal data read from mock file system: %v", err)
	}

	if !reflect.DeepEqual(readData, testTasksData) {
		t.Errorf("Saved data does not match the original data. Got: %v, Want: %v", readData, testTasksData)
	}
}
