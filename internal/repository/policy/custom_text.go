package policy

import (
	"encoding/json"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/asset"
)

const (
	customTextPrefix                   = "customtext."
	CustomTextSetEventType             = customTextPrefix + "set"
	CustomTextRemovedEventType         = customTextPrefix + "removed"
	CustomTextTemplateRemovedEventType = customTextPrefix + "template.removed"
	CustomTextFileUploadedEventType    = customTextPrefix + "file.uplaoded"
	CustomTextFileRemovedEventType     = customTextPrefix + "file.removed"
)

type CustomTextSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Key      string       `json:"key,omitempty"`
	Language language.Tag `json:"language,omitempty"`
	Text     string       `json:"text,omitempty"`
}

func (e *CustomTextSetEvent) Data() interface{} {
	return e
}

func (e *CustomTextSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextSetEvent(
	base *eventstore.BaseEvent,
	template,
	key,
	text string,
	language language.Tag,
) *CustomTextSetEvent {
	return &CustomTextSetEvent{
		BaseEvent: *base,
		Template:  template,
		Key:       key,
		Language:  language,
		Text:      text,
	}
}

func CustomTextSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &CustomTextSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-28dwe", "unable to unmarshal custom text")
	}

	return e, nil
}

type CustomTextRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Key      string       `json:"key,omitempty"`
	Language language.Tag `json:"language,omitempty"`
}

func (e *CustomTextRemovedEvent) Data() interface{} {
	return e
}

func (e *CustomTextRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextRemovedEvent(base *eventstore.BaseEvent, template, key string, language language.Tag) *CustomTextRemovedEvent {
	return &CustomTextRemovedEvent{
		BaseEvent: *base,
		Template:  template,
		Key:       key,
		Language:  language,
	}
}

func CustomTextRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &CustomTextRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-28sMf", "unable to unmarshal custom text removed")
	}

	return e, nil
}

type CustomTextTemplateRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Language language.Tag `json:"language,omitempty"`
}

func (e *CustomTextTemplateRemovedEvent) Data() interface{} {
	return e
}

func (e *CustomTextTemplateRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextTemplateRemovedEvent(base *eventstore.BaseEvent, template string, language language.Tag) *CustomTextTemplateRemovedEvent {
	return &CustomTextTemplateRemovedEvent{
		BaseEvent: *base,
		Template:  template,
		Language:  language,
	}
}

func CustomTextTemplateRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &CustomTextTemplateRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-mKKRs", "unable to unmarshal custom text message removed")
	}

	return e, nil
}

type CustomTextFileUploadedEvent struct {
	asset.AddedEvent
}

func (e *CustomTextFileUploadedEvent) Data() interface{} {
	return e
}

func (e *CustomTextFileUploadedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextFileUploadedEvent(base *eventstore.BaseEvent, storageKey string) *CustomTextFileUploadedEvent {
	return &CustomTextFileUploadedEvent{
		*asset.NewAddedEvent(base, storageKey),
	}
}

func CustomTextFileUploadedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &CustomTextFileUploadedEvent{*e.(*asset.AddedEvent)}, nil
}

type CustomTextFileRemovedEvent struct {
	asset.RemovedEvent
}

func (e *CustomTextFileRemovedEvent) Data() interface{} {
	return e
}

func (e *CustomTextFileRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextFileRemovedEvent(base *eventstore.BaseEvent, storageKey string) *CustomTextFileRemovedEvent {
	return &CustomTextFileRemovedEvent{
		*asset.NewRemovedEvent(base, storageKey),
	}
}

func CustomTextFileRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &CustomTextFileRemovedEvent{*e.(*asset.RemovedEvent)}, nil
}
