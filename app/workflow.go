package workmachine

import (
	"fmt"
	"reflect"
)

// Interface for serving crowdsourcing jobs.
type Assigner interface {
	// Publish a Batch to be done by workers.
	Execute(jobs chan Job, j Job)
}

type InputField string
type OutputField string

/*
 * Build a Job from a Task.
 */

type JobField struct {
	Id          string
	Value       string
	Description string
	Type        string
}

type Job struct {
	Task         Task
	InputFields  []JobField
	OutputFields []JobField
}

/*
 * Task is a way to define the Job that needs to be run.
 */
// crowd_type: output_image, output_text, input_text
type Task struct {
	Title       string
	Description string
	Price       int // In cents
	Tasks       interface{}
}

/*
 * Batch is a group of Jobs
 */

type Batch struct {
	Jobs    []Job
	Results []Job
}

func (b *Batch) Run(ss Assigner) {
	jobs := make(chan Job)

	for _, j := range b.Jobs {
		go ss.Execute(jobs, j)
	}

	for i := 0; i < len(b.Jobs); i++ {
		b.Results = append(b.Results, <-jobs)
	}

	fmt.Println(b.Results)
}

func NewBatch(task Task) (batch *Batch) {
	// Handle more of the task cases.
	tasks := task.Tasks

	batch = &Batch{}

	if reflect.TypeOf(tasks).Kind() != reflect.Slice {
		fmt.Println("Wtf kind of shit is this?")
		return nil
	}

	// Iterate over the Tasks
	s := reflect.ValueOf(tasks)
	for i := 0; i < s.Len(); i++ {
		// Figure out the information for one task.
		task := s.Index(i)

		job := Job{}

		// TODO NEED TO ADD TASK: Task: task
		// fmt.Println(task.Type())

		// Iterate over the fields of a struct
		for j := 0; j < task.NumField(); j++ {

			switch task.Type().Field(j).Type.Name() {
			case "InputField":
				job.InputFields = append(job.InputFields, JobField{
					// fmt.Println(task.Type().Field(j).Name)
					Id:          task.Type().Field(j).Tag.Get("crowd_id"),
					Description: task.Type().Field(j).Tag.Get("crowd_desc"),
					Type:        task.Type().Field(j).Tag.Get("crowd_type"),
					Value:       task.Field(j).String(),
				})
			case "OutputField":
				job.OutputFields = append(job.OutputFields, JobField{
					// fmt.Println(task.Type().Field(j).Name)
					Id:          task.Type().Field(j).Tag.Get("crowd_id"),
					Description: task.Type().Field(j).Tag.Get("crowd_desc"),
					Type:        task.Type().Field(j).Tag.Get("crowd_type"),
				})

			default:
				// fmt.Println("Unknown")
			}
		}

		batch.Jobs = append(batch.Jobs, job)
	}

	return
}