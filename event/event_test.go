package event

import (
	"testing"
	"time"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct{}
type TestEventHandler2 struct{}
type TestEventHandler3 struct{}

func (h *TestEventHandler) Handle(event EventInterface)  {}
func (h *TestEventHandler2) Handle(event EventInterface) {}
func (h *TestEventHandler3) Handle(event EventInterface) {}

type EventDispatcherTestSuite struct {
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler2
	handler3        TestEventHandler3
	eventDispatcher *EventDispatcher
}

func TestEventDispatcher(t *testing.T) {
	suite := &EventDispatcherTestSuite{}
	suite.SetupTest()

	t.Run("TestEventDispatcher_Register", suite.TestEventDispatcher_Register)
	// clear the event dispatcher after the test
	suite.eventDispatcher.Clear()
	t.Run("TestEventDispatcher_Register_WithSameHandler", suite.TestEventDispatcher_Register_WithSameHandler)
	// clear the event dispatcher after the test
	suite.eventDispatcher.Clear()
	t.Run("TestEventDispatcher_Clear", suite.TestEventDispatcher_Clear)
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
	suite.handler = TestEventHandler{}
	suite.handler2 = TestEventHandler2{}
	suite.handler3 = TestEventHandler3{}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register(t *testing.T) {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(suite.eventDispatcher.handlers[suite.event.GetName()]) != 1 {
		t.Errorf("expected 1 handler, got %d", len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	}
	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(suite.eventDispatcher.handlers[suite.event.GetName()]) != 2 {
		t.Errorf("expected 2 handlers, got %d", len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	}
	if suite.eventDispatcher.handlers[suite.event.GetName()][0] != &suite.handler {
		t.Errorf("expected handler1, got %v", suite.eventDispatcher.handlers[suite.event.GetName()][0])
	}
	if suite.eventDispatcher.handlers[suite.event.GetName()][1] != &suite.handler2 {
		t.Errorf("expected handler2, got %v", suite.eventDispatcher.handlers[suite.event.GetName()][1])
	}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler(t *testing.T) {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(suite.eventDispatcher.handlers[suite.event.GetName()]) != 1 {
		t.Errorf("expected 1 handler, got %d", len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	}
	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	if err != ErrHandlerAlreadyRegistered {
		t.Errorf("expected error %v, got %v", ErrHandlerAlreadyRegistered, err)
	}
	if len(suite.eventDispatcher.handlers[suite.event.GetName()]) != 1 {
		t.Errorf("expected 1 handler, got %d", len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	}

}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear(t *testing.T) {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(suite.eventDispatcher.handlers[suite.event.GetName()]) != 1 {
		t.Errorf("expected 1 handler, got %d", len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	}
	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(suite.eventDispatcher.handlers[suite.event.GetName()]) != 2 {
		t.Errorf("expected 2 handlers, got %d", len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	}
	suite.eventDispatcher.Clear()
	if len(suite.eventDispatcher.handlers) != 0 {
		t.Errorf("expected 0 handlers, got %d", len(suite.eventDispatcher.handlers))
	}
}
