package jira

import (
	"errors"
	"github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"tix/ticket"
)

var issue = &jira.Issue{
	Fields: &jira.IssueFields{
		Summary: "An issue",
	},
	ID:  "1",
	Key: "key",
}
var jiraFields = []jira.Field{
	{ID: "field1", Name: "epic name"},
	{ID: "field2", Name: "epic link"},
}

func TestCreator_CreateTickets_NoFieldsError(t *testing.T) {
	api := NewMockApi(nil)
	creator := NewCreatorWithEpics(api)
	rootTicket1 := ticket.NewTicket()
	rootTicket2 := ticket.NewTicket()
	childTicket := ticket.NewTicket()
	rootTicket1.AddSubticket(childTicket)
	tickets := []*ticket.Ticket{rootTicket1, rootTicket2}

	creator.CreateTickets(tickets)
	api.AssertCalled(t, "GetIssueFieldList")
	api.AssertCalled(t, "CreateIssue", mock.AnythingOfType("*jira.Issue"))
	api.AssertNumberOfCalls(t, "CreateIssue", 3)
}

func TestCreator_CreateTickets_WithFieldsError(t *testing.T) {
	api := NewMockApi(errors.New("oh noes"))
	creator := NewCreatorWithEpics(api)
	rootTicket1 := ticket.NewTicket()
	rootTicket2 := ticket.NewTicket()
	childTicket := ticket.NewTicket()
	rootTicket1.AddSubticket(childTicket)
	tickets := []*ticket.Ticket{rootTicket1, rootTicket2}

	creator.CreateTickets(tickets)
	api.AssertCalled(t, "GetIssueFieldList")
	api.AssertNotCalled(t, "CreateIssue", mock.AnythingOfType("*jira.Issue"))
}

func TestCreator_UpdateTickets(t *testing.T) {
	api := NewMockApi(nil)
	creator := NewCreatorWithEpics(api)
	rootTicket1 := ticket.NewTicketWithFields(map[string]interface{}{"update_ticket": "JIRA123"})
	rootTicket2 := ticket.NewTicket()
	childTicket := ticket.NewTicket()
	rootTicket1.AddSubticket(childTicket)
	tickets := []*ticket.Ticket{rootTicket1, rootTicket2}

	creator.CreateTickets(tickets)
	api.AssertCalled(t, "GetIssueFieldList")
	api.AssertCalled(t, "CreateIssue", mock.AnythingOfType("*jira.Issue"))
	api.AssertCalled(t, "UpdateIssue", mock.AnythingOfType("*jira.Issue"))
	api.AssertNumberOfCalls(t, "CreateIssue", 2)
	api.AssertNumberOfCalls(t, "UpdateIssue", 1)
}

func TestNewCreatorWithEpics(t *testing.T) {
	api := NewMockApi(nil)
	creator := NewCreatorWithEpics(api)
	assert.Equal(t, creator.startingTicketLevel, 0)
}

func TestNewCreatorWithoutEpics(t *testing.T) {
	api := NewMockApi(nil)
	creator := NewCreatorWithoutEpics(api)
	assert.Equal(t, creator.startingTicketLevel, 1)
}

type MockApi struct {
	mock.Mock
}

func NewMockApi(fieldsError error) *MockApi {
	api := &MockApi{}
	api.On("GetIssueFieldList").Return(jiraFields, fieldsError)
	api.On("CreateIssue", mock.AnythingOfType("*jira.Issue")).Return(issue, nil)
	api.On("UpdateIssue", mock.AnythingOfType("*jira.Issue")).Return(issue, nil)
	return api
}

func (m *MockApi) CreateIssue(issue *jira.Issue) (*jira.Issue, error) {
	args := m.Called(issue)
	return args.Get(0).(*jira.Issue), args.Error(1)
}

func (m *MockApi) GetIssueFieldList() ([]jira.Field, error) {
	args := m.Called()
	return args.Get(0).([]jira.Field), args.Error(1)
}

func (m *MockApi) UpdateIssue(issue *jira.Issue) (*jira.Issue, error) {
	args := m.Called(issue)
	return args.Get(0).(*jira.Issue), args.Error(1)
}
