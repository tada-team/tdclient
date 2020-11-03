package tdclient

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tada-team/tdproto"
	"github.com/tada-team/tdproto/tdapi"
)

func (s Session) Ping() error {
	resp := new(struct {
		tdapi.Resp
		Result string `json:"result"`
	})
	return s.doGet("/api/v4/ping", nil, resp)
}

func (s Session) Me(teamUid string) (tdproto.Contact, error) {
	resp := new(struct {
		tdapi.Resp
		Result tdproto.Team `json:"result"`
	})

	if !tdproto.ValidUid(teamUid) {
		return tdproto.Contact{}, errors.New("invalid team uid")
	}

	if err := s.doGet("/api/v4/teams/"+teamUid, nil, resp); err != nil {
		return tdproto.Contact{}, err
	}

	if !resp.Ok {
		return tdproto.Contact{}, resp.Error
	}

	return resp.Result.Me, nil
}

func (s Session) Contacts(teamUid string) ([]tdproto.Contact, error) {
	resp := new(struct {
		tdapi.Resp
		Result []tdproto.Contact `json:"result"`
	})

	if !tdproto.ValidUid(teamUid) {
		return resp.Result, errors.New("invalid team uid")
	}

	if err := s.doGet(fmt.Sprintf("/api/v4/teams/%s/contacts/", teamUid), nil, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) AddContact(teamUid string, phone string) (tdproto.Contact, error) {
	req := map[string]interface{}{
		"phone": phone,
	}

	resp := new(struct {
		tdapi.Resp
		Result tdproto.Contact `json:"result"`
	})

	if err := s.doPost(fmt.Sprintf("/api/v4/teams/%s/contacts", teamUid), req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) AuthBySmsSendCode(phone string) (tdapi.SmsCode, error) {
	req := map[string]interface{}{
		"phone": phone,
	}

	resp := new(struct {
		tdapi.Resp
		Result tdapi.SmsCode `json:"result"`
	})

	if err := s.doPost("/api/v4/auth/sms/send-code", req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) AuthBySmsGetToken(phone, code string) (tdapi.Auth, error) {
	req := map[string]interface{}{
		"phone": phone,
		"code":  code,
	}

	resp := new(struct {
		tdapi.Resp
		Result tdapi.Auth `json:"result"`
	})

	if err := s.doPost("/api/v4/auth/sms/get-token", req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) AuthByPasswordGetToken(username, password string) (tdapi.Auth, error) {
	req := map[string]string{
		"username": username,
		"password": password,
	}

	resp := new(struct {
		tdapi.Resp
		Result tdapi.Auth `json:"result"`
	})

	if err := s.doPost("/api/v4/auth/password/get-token", req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) SendPlaintextMessage(teamUid string, chat tdproto.JID, text string) (tdproto.Message, error) {
	req := new(tdapi.Message)
	req.Type = tdproto.MediatypePlain
	req.Text = text

	req.MessageUid = uuid.New().String()

	resp := new(struct {
		tdapi.Resp
		Result tdproto.Message `json:"result"`
	})

	if err := s.doPost(fmt.Sprintf("/api/v4/teams/%s/chats/%s/messages", teamUid, chat), req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) GetMessages(teamUid string, chat tdproto.JID) (tdproto.ChatMessages, error) {
	resp := new(struct {
		tdapi.Resp
		Result tdproto.ChatMessages `json:"result"`
	})

	if err := s.doGet(fmt.Sprintf("/api/v4/teams/%s/messages/%s", teamUid, chat), resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) DeleteMessage(teamUid string, chat tdproto.JID, MsgId string) (tdproto.ChatMessages, error) {

	resp := new(struct {
		tdapi.Resp
		Result tdproto.ChatMessages `json:"result"`
	})

	if err := s.doDelete(fmt.Sprintf("/api/v4/teams/%s/groups/%s/message/%s", teamUid, chat, MsgId), resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) GetMessagesOldMsg(teamUid string, chat tdproto.JID, lastMsgId string) (tdproto.ChatMessages, error) {

	param := map[string]string{}
	param["old_from"] = lastMsgId

	resp := new(struct {
		tdapi.Resp
		Result tdproto.ChatMessages `json:"result"`
	})

	if err := s.doGetParams(fmt.Sprintf("/api/v4/teams/%s/messages/%s", teamUid, chat), resp, param); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) CreateTask(teamUid string, req tdapi.Task) (tdproto.Chat, error) {
	resp := new(struct {
		tdapi.Resp
		Result tdproto.Chat `json:"result"`
	})

	if err := s.doPost(fmt.Sprintf("/api/v4/teams/%s/tasks", teamUid), req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) CreateGroup(teamUid string, req tdapi.Group) (tdproto.Chat, error) {
	resp := new(struct {
		tdapi.Resp
		Result tdproto.Chat `json:"result"`
	})

	if err := s.doPost(fmt.Sprintf("/api/v4/teams/%s/groups", teamUid), req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) GetGroups(teamUid string) ([]tdproto.Chat, error) {
	resp := new(struct {
		tdapi.Resp
		Result []tdproto.Chat `json:"result"`
	})

	if err := s.doGet(fmt.Sprintf("/api/v4/teams/%s/groups", teamUid), nil, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) AddGroupMember(teamUid string, group, contact tdproto.JID) (tdproto.GroupMembership, error) {
	req := tdapi.GroupMember{
		Jid:    contact,
		Status: tdproto.GroupMember,
	}

	resp := new(struct {
		tdapi.Resp
		Result tdproto.GroupMembership `json:"result"`
	})

	if err := s.doPost(fmt.Sprintf("/api/v4/teams/%s/groups/%s/members", teamUid, group), req, resp); err != nil {
		return resp.Result, err
	}

	if !resp.Ok {
		return resp.Result, resp.Error
	}

	return resp.Result, nil
}

func (s Session) GroupMembers(teamUid string, group tdproto.JID) ([]tdproto.GroupMembership, error) {
	type MembersParams struct {
		Members []tdproto.GroupMembership `json:"members"`
	}
	resp := new(struct {
		tdapi.Resp
		Result MembersParams `json:"result"`
	})

	if !tdproto.ValidUid(teamUid) {
		return resp.Result.Members, errors.New("invalid team uid")
	}

	if err := s.doGet(fmt.Sprintf("/api/v4/teams/%s/groups/%s/members", teamUid, group), nil, resp); err != nil {
		return resp.Result.Members, err
	}

	if !resp.Ok {
		return resp.Result.Members, resp.Error
	}

	return resp.Result.Members, nil
}

func (s Session) DropGroupMember(teamUid string, group, contact tdproto.JID) error {
	resp := new(tdapi.Resp)

	if !tdproto.ValidUid(teamUid) {
		return InvalidTeamUid
	}

	if err := s.doDelete(fmt.Sprintf("/api/v4/teams/%s/groups/%s/members/%s", teamUid, group, contact), resp); err != nil {
		return err
	}

	if !resp.Ok {
		return resp.Error
	}

	return nil
}

func (s Session) DropGroup(teamUid string, group tdproto.JID) error {
	resp := new(tdapi.Resp)

	if !tdproto.ValidUid(teamUid) {
		return InvalidTeamUid
	}

	if err := s.doDelete(fmt.Sprintf("/api/v4/teams/%s/groups/%s", teamUid, group), resp); err != nil {
		return err
	}

	if !resp.Ok {
		return resp.Error
	}

	return nil
}

func (s Session) GetChats(teamUid string, f *tdapi.ChatFilter) ([]tdproto.Chat, error) {
	resp := new(struct {
		tdapi.Resp
		Result tdproto.PaginatedChats `json:"result"`
	})

	if f == nil {
		f = new(tdapi.ChatFilter)
	}

	if f.Limit == 0 {
		f.Limit = 100
	}

	result := make([]tdproto.Chat, 0)
	for {
		if err := s.doGet(fmt.Sprintf("/api/v4/teams/%s/chats", teamUid), f, resp); err != nil {
			return result, err
		}

		if !resp.Ok {
			return result, resp.Error
		}

		if len(resp.Result.Objects) == 0 {
			break
		}

		f.Offset += f.Limit
		result = append(result, resp.Result.Objects...)
	}

	return result, nil
}
