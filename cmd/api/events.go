package main

import (
	"api-in-gin/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createEvent creates an event and inserts it into database
func (app *application) createEvent(c *gin.Context) {
	var event database.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// getAllEvents returns the list of all events
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// getEvent returns an event by ID
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
	}

	event, err := app.models.Events.Get(id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
	}

	c.JSON(http.StatusOK, event)
}

// updateEvent updates an event by ID
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	updatedEvent := database.Event{}
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.ID = id
	if err = app.models.Events.Update(&updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// deleteEvent deletes an event by id
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event"})
	}

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	event, err := app.models.Events.Get(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.ID, userToAdd.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventID: event.ID,
		UserID:  userToAdd.ID,
	}
	err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	attendees, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendees for event"})
		return
	}
	if attendees == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attendees not found"})
		return
	}

	c.JSON(http.StatusOK, attendees)
}

func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	attendeeID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attendee id"})
		return
	}

	event, err := app.models.Events.Get(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	attendee, err := app.models.Attendees.GetByEventAndAttendee(event.ID, attendeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendee"})
		return
	}
	if attendee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attendee not found"})
		return
	}

	err = app.models.Attendees.Delete(eventID, attendeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete attendee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attendee id"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve events"})
		return
	}
	if events == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "events not found"})
		return
	}

	c.JSON(http.StatusOK, events)
}
