package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Ajesh8/UserTasks/pkg/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

type Task struct {
	gorm.Model
	Title       string       `gorm:"" json:"title"`
	Description string       `json:"description"`
	DueDate     time.Time    `json:"dueDate"`
	Status      string       `json:"status"`
	Subtask     []Task       `gorm:"ForeignKey:ParentID" json:"subTasks"`
	ParentID    *uint        `json:"parentID"`
	UserID      uint         `json:userId`
	EmailAlerts []EmailAlert `json:"emailAlerts"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&EmailAlert{}, &User{}, &Task{})
}

func (b *Task) CreateTask(user *User) (*Task, error) {
	if !validateTask(*b) {
		return nil, errors.New("task/subtasks' mandatory fields not found")
	}
	b.Status = "Pending"
	b.UserID = user.ID
	for i := range b.Subtask {
		b.Subtask[i].Status = "Pending"
		b.Subtask[i].UserID = user.ID
	}
	db.Create(&b)
	return b, nil
}

func validateTask(task Task) bool {
	if task.Title == "" {
		return false
	}
	for _, subtask := range task.Subtask {
		if !validateTask(subtask) {
			return false
		}
	}
	return true
}

func FindTaskById(taskId int64) *Task {
	task := &Task{}
	db.Preload(clause.Associations).Where("ID=?", taskId).Find(&task)
	return task
}

func UpdateTaskStatus(taskId int64, status string) *Task {
	task := FindTaskById(taskId)
	task.Status = status
	if status == "Completed" {
		for i := range task.Subtask {
			task.Subtask[i].Status = status
			db.Save(task.Subtask[i])
		}
	}
	db.Save(task)
	return task
}

func GetAllTask() []Task {
	var Tasks []Task
	db.Preload(clause.Associations).Order("due_date").Find(&Tasks)
	return Tasks
}

func SearchTask(searchText string) []Task {
	var Tasks []Task
	db.Where("title like ?", "%"+searchText+"%").Order("due_date").Find(&Tasks)
	return Tasks
}

func FilterTask(dueText string) []Task {
	var Tasks []Task
	db.Order("due_date").Find(&Tasks)
	t := time.Now()
	fmt.Println(dueText)
	year, month, day := t.Date()
	todayMidnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
	fmt.Println(todayMidnight)
	weekendMidnight := todayMidnight.AddDate(0, 0, 7-int(t.Weekday())%7)
	fmt.Println(weekendMidnight)
	var Result []Task

	for _, task := range Tasks {
		if dueText == "Overdue" {
			if task.DueDate.Before(time.Now()) {
				Result = append(Result, task)
			}
		} else {
			if dueText == "Today" {
				if DateEqual(task.DueDate, t) {
					Result = append(Result, task)
				}
			} else if dueText == "This Week" {
				if task.DueDate.After(weekendMidnight.AddDate(0, 0, -7)) && task.DueDate.Before(weekendMidnight) {
					Result = append(Result, task)
				}
			} else if dueText == "Next Week" {
				if task.DueDate.After(weekendMidnight) && task.DueDate.Before(weekendMidnight.AddDate(0, 0, 7)) {
					Result = append(Result, task)
				}
			} else {
				fmt.Print("No proper due term found. Finding tasks due today")
				if DateEqual(task.DueDate, t) {
					Result = append(Result, task)
				}
			}

		}
	}
	return Result
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
