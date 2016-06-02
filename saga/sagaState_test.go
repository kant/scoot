package saga

import "bytes"
import "fmt"
import "testing"

func TestsagaStateFactory(t *testing.T) {

	sagaId := "testSaga"
	job := []byte{0, 1, 2, 3, 4, 5}

	state, _ := sagaStateFactory("testSaga", job)
	if state.sagaId != sagaId {
		t.Error(fmt.Sprintf("SagaState SagaId should be the same as the SagaId passed to Factory Method"))
	}

	if !bytes.Equal(state.Job(), job) {
		t.Error(fmt.Sprintf("SagaState Job should be the same as the supplied Job passed to Factory Method"))
	}
}

func TestSagaState_AbortSaga(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	if state.IsSagaAborted() {
		t.Error("IsSagaAborted should return false")
	}

	state, err := updateSagaState(state, AbortSagaMessageFactory(sagaId))
	if err != nil {
		t.Error(fmt.Sprintf("AbortSaga Failed Unexpected %s", err))
	}

	if !state.IsSagaAborted() {
		t.Error("IsSagaAborted should return true")
	}
}

func TestSagaState_StartTask(t *testing.T) {
	sagaId := "testSaga"
	taskId := "task1"
	state, _ := sagaStateFactory(sagaId, nil)

	if state.IsTaskStarted(taskId) {
		t.Error("TaskStarted should return false")
	}

	state, err := updateSagaState(state, StartTaskMessageFactory(sagaId, taskId))
	if err != nil {
		t.Error(fmt.Sprintf("StartTask Failed Unexpected %s", err))
	}

	if !state.IsTaskStarted(taskId) {
		t.Error("TaskStarted should return true")
	}
}

func TestSagaState_EndTask(t *testing.T) {
	sagaId := "testSaga"
	taskId := "task1"
	state, _ := sagaStateFactory(sagaId, nil)

	if state.IsTaskCompleted(taskId) {
		t.Error("TaskCompleted should return false")
	}

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, taskId),
		EndTaskMessageFactory(sagaId, taskId, nil),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType.String(), err))
		}
	}

	if !state.IsTaskCompleted(taskId) {
		t.Error("TaskCompleted should return true")
	}
}

func TestSagaState_EndTaskBeforeStartTaskFails(t *testing.T) {
	sagaId := "testSaga"
	taskId := "task1"
	state, _ := sagaStateFactory(sagaId, nil)

	var err error
	state, err = updateSagaState(state, EndTaskMessageFactory(sagaId, taskId, nil))
	if err == nil {
		t.Error("EndTask Should Fail When Written Before Start Task")
	}
}

func TestSagaState_EndSaga(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	if state.IsSagaCompleted() {
		t.Error("IsSagaCompleted should return false")
	}

	var err error
	state, err = updateSagaState(state, EndSagaMessageFactory(sagaId))
	if err != nil {
		t.Error(fmt.Sprintf("EndSaga Failed Unexpected %s", err))
	}

	if !state.IsSagaCompleted() {
		t.Error("IsSagaCompleted should return true")
	}
}

func TestSagaState_EndSagaBeforeAllTasksCompleted(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task1"),
		StartTaskMessageFactory(sagaId, "task2"),
		StartTaskMessageFactory(sagaId, "task3"),
		EndTaskMessageFactory(sagaId, "task2", nil),
		EndTaskMessageFactory(sagaId, "task1", nil),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType.String(), err))
		}
	}

	var err error
	state, err = updateSagaState(state, EndSagaMessageFactory(sagaId))
	if err == nil {
		t.Error("EndSaga Should Fail when not all tasks completed")
	}
}

func TestSagaState_EndSagaBeforeAllCompTasksCompleted(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task1"),
		AbortSagaMessageFactory(sagaId),
		StartCompTaskMessageFactory(sagaId, "task1"),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	var err error
	state, err = updateSagaState(state, EndSagaMessageFactory(sagaId))
	if err == nil {
		t.Error("EndSaga Should Fail when not all comp tasks completed")
	}
}

func TestSagaState_StartCompTask(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)
	taskId := "task1"

	if state.IsCompTaskStarted(taskId) {
		t.Error("IsCompTaskStarted should return false")
	}

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, taskId),
		AbortSagaMessageFactory(sagaId),
		StartCompTaskMessageFactory(sagaId, taskId),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	if !state.IsCompTaskStarted(taskId) {
		t.Error("IsCompTaskStarted should return true")
	}
}

func TestSagaState_StartCompTaskNoStartTask(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task1"),
		AbortSagaMessageFactory(sagaId),
		StartCompTaskMessageFactory(sagaId, "task1"),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	var err error
	state, err = updateSagaState(state, StartCompTaskMessageFactory(sagaId, "task2"))
	if err == nil {
		t.Error("StartCompTask Should Fail when not all comp tasks completed")
	}
}

func TestSagaState_StartCompTaskNoAbort(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task1"),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	var err error
	state, err = updateSagaState(state, StartCompTaskMessageFactory(sagaId, "task1"))
	if err == nil {
		t.Error("EndSaga Should Fail when not all comp tasks completed")
	}
}

func TestSagaState_EndCompTask(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)
	taskId := "task1"

	if state.IsCompTaskCompleted(taskId) {
		t.Error("IsCompTaskCompleted should return false")
	}

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, taskId),
		AbortSagaMessageFactory(sagaId),
		StartCompTaskMessageFactory(sagaId, taskId),
		EndCompTaskMessageFactory(sagaId, taskId, nil),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	if !state.IsCompTaskCompleted(taskId) {
		t.Error("IsCompTaskCompleted should return true")
	}
}

func TestSagaState_EndCompTaskNoStartTask(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		AbortSagaMessageFactory(sagaId),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	var err error
	state, err = updateSagaState(state, EndCompTaskMessageFactory(sagaId, "task2", nil))
	if err == nil {
		t.Error("StartCompTask Should Fail when not all comp tasks completed")
	}
}

func TestSagaState_EndCompTaskNoStartCompTask(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task2"),
		AbortSagaMessageFactory(sagaId),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	var err error
	state, err = updateSagaState(state, EndCompTaskMessageFactory(sagaId, "task2", nil))
	if err == nil {
		t.Error("StartCompTask Should Fail when not all comp tasks completed")
	}
}

func TestSagaState_EndCompTaskNoAbort(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task2"),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	var err error
	state, err = updateSagaState(state, EndCompTaskMessageFactory(sagaId, "task2", nil))
	if err == nil {
		t.Error("StartCompTask Should Fail when not all comp tasks completed")
	}
}

func TestSagaState_SuccessfulSaga(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task1"),
		StartTaskMessageFactory(sagaId, "task2"),
		StartTaskMessageFactory(sagaId, "task3"),
		EndTaskMessageFactory(sagaId, "task2", nil),
		EndTaskMessageFactory(sagaId, "task3", nil),
		EndTaskMessageFactory(sagaId, "task1", nil),
		EndSagaMessageFactory(sagaId),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	if !state.IsSagaCompleted() {
		t.Error("Expected Saga to be Completed")
	}
}

func TestSagaState_AbortedSaga(t *testing.T) {
	sagaId := "testSaga"
	state, _ := sagaStateFactory(sagaId, nil)

	msgs := []sagaMessage{
		StartTaskMessageFactory(sagaId, "task1"),
		StartTaskMessageFactory(sagaId, "task2"),
		StartTaskMessageFactory(sagaId, "task3"),
		EndTaskMessageFactory(sagaId, "task2", nil),
		AbortSagaMessageFactory(sagaId),
		StartCompTaskMessageFactory(sagaId, "task1"),
		StartCompTaskMessageFactory(sagaId, "task2"),
		EndCompTaskMessageFactory(sagaId, "task2", nil),
		EndCompTaskMessageFactory(sagaId, "task1", nil),
		StartCompTaskMessageFactory(sagaId, "task3"),
		EndCompTaskMessageFactory(sagaId, "task3", nil),
		EndSagaMessageFactory(sagaId),
	}

	for _, msg := range msgs {
		var err error
		state, err = updateSagaState(state, msg)
		if err != nil {
			t.Error(fmt.Sprintf("Applying Saga Message %s Failed Unexpectedly %s", msg.msgType, err))
		}
	}

	if !state.IsSagaCompleted() {
		t.Error("Expected Saga to be Completed")
	}
}

func TestSagaState_ValidateSagaId(t *testing.T) {
	err := validateSagaId("")
	if err == nil {
		t.Error(fmt.Sprintf("Invalid Saga Id Should Return Error"))
	}
}

func TestSagaState_ValidateTaskId(t *testing.T) {
	err := validateTaskId("")
	if err == nil {
		t.Error(fmt.Sprintf("Invalid Task Id Should Return Error"))
	}
}
