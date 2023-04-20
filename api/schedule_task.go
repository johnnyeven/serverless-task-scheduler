package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"reflect"
	"severless-task-scheduler/db/model"
	"strconv"
)

type PredictRequest interface {
	Parse(r io.Reader, task *model.Task) error
	Json() []byte
}

type GradioRequest struct {
	Data []any `json:"data"`
}

func (g *GradioRequest) Json() []byte {
	marshal, _ := json.Marshal(g)
	return marshal
}

func (g *GradioRequest) Parse(r io.Reader, task *model.Task) error {
	all, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	content := make(map[string]any)
	err = json.Unmarshal(all, &content)
	if err != nil {
		return err
	}
	g.Data = append(g.Data, task.ID)
	g.Data = append(g.Data, content["prompt"])
	g.Data = append(g.Data, DefaultNegativePrompt)
	g.Data = append(g.Data, DefaultNumInferenceSteps)
	g.Data = append(g.Data, DefaultWidth)
	g.Data = append(g.Data, DefaultHeight)
	g.Data = append(g.Data, DefaultGuidanceScale)
	g.Data = append(g.Data, DefaultRandSeed)
	return nil
}

type Model struct {
	Name string `json:"name"`
	Api  string `json:"api"`
}

type TaskParameter struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

const (
	DefaultNegativePrompt    = "nsfw,(worst quality:2),(low quality:2)"
	DefaultNumInferenceSteps = 20
	DefaultWidth             = 512
	DefaultHeight            = 512
	DefaultGuidanceScale     = 7
	DefaultRandSeed          = -1
)

var modelRequest = map[string]reflect.Type{
	"openjourney":    reflect.TypeOf(GradioRequest{}),
	"anything":       reflect.TypeOf(GradioRequest{}),
	"waifu":          reflect.TypeOf(GradioRequest{}),
	"nerverendDream": reflect.TypeOf(GradioRequest{}),
}

var models map[string]Model

var client = &http.Client{
	Timeout: 0,
}

func init() {
	models = make(map[string]Model)
	configJson := os.Getenv("MODEL_CONFIG")
	if configJson == "" {
		panic("MODEL_CONFIG is empty")
	}
	err := json.Unmarshal([]byte(configJson), &models)
	if err != nil {
		panic(err)
	}
}

func ScheduleTask(w http.ResponseWriter, r *http.Request) {
	// 查询所有状态为1的任务
	limitConfig := os.Getenv("SCHEDULE_TASK_LIMIT")
	if limitConfig == "" {
		logrus.Error("SCHEDULE_TASK_LIMIT is empty")
		responseError(w, errors.New("SCHEDULE_TASK_LIMIT is empty"))
		return
	}
	limit, err := strconv.Atoi(limitConfig)
	if err != nil {
		logrus.Error("SCHEDULE_TASK_LIMIT is not a number")
		responseError(w, err)
		return
	}
	tasks, err := query.Task.Where(query.Task.Status.Eq(1)).Order(query.Task.ID).Limit(limit).Find()
	if err != nil {
		logrus.Errorf("query task error: %v", err)
		responseError(w, err)
		return
	}
	for _, task := range tasks {
		// 生成任务参数
		taskParameter := TaskParameter{}
		err = json.Unmarshal([]byte(task.Parameter), &taskParameter)
		if err != nil {
			// 更新任务状态为失败
			_, _ = query.Task.Where(query.Task.ID.Eq(task.ID)).UpdateColumn(query.Task.Status, int32(Fail))
			continue
		}
		m, ok := models[taskParameter.Model]
		if !ok {
			// 更新任务状态为失败
			_, _ = query.Task.Where(query.Task.ID.Eq(task.ID)).UpdateColumn(query.Task.Status, int32(Fail))
			continue
		}
		// 调用模型API
		go call(m, task)
	}
	responseEmpty(w)
}

func call(m Model, task *model.Task) {
	requestReflect, ok := modelRequest[m.Name]
	if !ok {
		// 更新任务状态为失败
		logrus.Errorf("model requestReflect %s not found", m.Name)
		_, _ = query.Task.Where(query.Task.ID.Eq(task.ID)).UpdateColumn(query.Task.Status, int32(Fail))
		return
	}
	request := reflect.New(requestReflect).Interface().(PredictRequest)
	err := request.Parse(bytes.NewReader([]byte(task.Parameter)), task)
	if err != nil {
		// 更新任务状态为失败
		logrus.Errorf("parse task parameter error: %v", err)
		_, _ = query.Task.Where(query.Task.ID.Eq(task.ID)).UpdateColumn(query.Task.Status, int32(Fail))
		return
	}

	// 更新任务状态为处理中
	_, _ = query.Task.Where(query.Task.ID.Eq(task.ID)).UpdateColumn(query.Task.Status, int32(Running))
	_, _ = client.Post(m.Api, "application/json", bytes.NewReader(request.Json()))
}
