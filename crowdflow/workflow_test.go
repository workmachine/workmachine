package crowdflow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Workflow", func() {

	var (
		tasks    []Research
		research TaskConfig
		batch    *Batch
	)

	BeforeEach(func() {
		tasks = []Research{
			Research{Name: "Michelangelo"},
			Research{Name: "Leonardo"},
		}

		research = TaskConfig{
			Title:       "Research the fields",
			Description: "Research the person.",
			Price:       1,
			Tasks:       tasks,
			Write:       func(j *MetaJob) {},
		}

		batch = NewBatch(research)
	})

	Describe("Batch", func() {
		Describe("TaskConfig", func() {
			It("has the correct TaskConfig", func() {
				Expect(batch.TaskConfig.Title).To(Equal("Research the fields"))
				Expect(batch.TaskConfig.Description).To(Equal("Research the person."))
				Expect(batch.TaskConfig.Price).To(Equal(uint(1)))
			})
		})

		Describe("MetaJobs", func() {
			It("has the correct MetaJobs", func() {
				Expect(len(batch.MetaJobs)).To(Equal(len(tasks)))
			})

			It("has the correct titles and descriptions", func() {
				Expect(len(batch.MetaJobs)).To(Equal(len(tasks)))

				for _, m := range batch.MetaJobs {
					Expect(m.TaskConfig.Title).To(Equal("Research the fields"))
					Expect(m.TaskConfig.Description).To(Equal("Research the person."))
				}
			})

			It("has the correct input fields", func() {
				m := batch.MetaJobs

				Expect(m[0].InputFields[0].Id).To(Equal("name"))
				Expect(m[0].InputFields[0].Value).To(Equal("Michelangelo"))
				Expect(m[0].InputFields[0].Description).To(Equal("The name of the person to research."))
				Expect(m[0].InputFields[0].Type).To(Equal(JobFieldType("")))

				Expect(m[1].InputFields[0].Id).To(Equal("name"))
				Expect(m[1].InputFields[0].Value).To(Equal("Leonardo"))
				Expect(m[1].InputFields[0].Description).To(Equal("The name of the person to research."))
				Expect(m[1].InputFields[0].Type).To(Equal(JobFieldType("")))
			})

			It("has the correct output fields", func() {
				for _, m := range batch.MetaJobs {
					Expect(m.OutputFields[0].Id).To(Equal("born"))
					Expect(m.OutputFields[1].Id).To(Equal("is_painter"))
					Expect(m.OutputFields[1].Value).To(Equal(""))
					Expect(m.OutputFields[1].Description).To(Equal("Was the person a painter?"))
					Expect(m.OutputFields[1].Type).To(Equal(JobFieldType("checkbox")))
				}

			})
		})

		Describe("Run", func() {
			BeforeEach(func() {
				go func() {
					batch.Run()
				}()
			})

			It("Has the correct number of assignments", func() {
				time.Sleep(1000) // Need time for it to run.

				Expect(len(AvailableAssignments)).To(Equal(len(batch.MetaJobs)))
				for _, a := range AvailableAssignments {
					Expect(a.Assigned).To(Equal(false))
					// Expect(a.Id).To(Equal
					// Expect(a.StartedAt)
					Expect(a.Finished).To(Equal(false))
				}

				a := AvailableAssignments[0]

				Expect(&a).To(Equal(AvailableAssignments.GetUnfinished()))

				a.SharedAssignment.Assign()
				Expect(a.SharedAssignment.Assigned).To(Equal(true))
				Expect(a.SharedAssignment.TryToAssign()).To(Equal(false))

				a2 := AvailableAssignments[1]
				Expect(a2.SharedAssignment.Assigned).To(Equal(false))
				Expect(a2.SharedAssignment.TryToAssign()).To(Equal(true))
			})

			PDescribe("Job", func() {

			})
		})

	})
})
