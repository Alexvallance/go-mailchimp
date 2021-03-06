package mailchimp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Alexvallance/go-mailchimp/v3/status"
)

// Unsubscribe ...
func (c *Client) UnSubscribe(listID string, email string, mergeFields map[string]interface{}) (*MemberResponse, error) {
	// Make request
	params := map[string]interface{}{
		"email_address": email,
		"status":        status.Unsubscribed,
		"merge_fields":  mergeFields,
	}
	resp, err := c.do(
		"PUT",
		fmt.Sprintf("/lists/%s/members/", listID),
		&params,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Allow any success status (2xx)
	if resp.StatusCode/100 == 2 {
		// Unmarshal response into MemberResponse struct
		memberResponse := new(MemberResponse)
		if err := json.Unmarshal(data, memberResponse); err != nil {
			return nil, err
		}
		return memberResponse, nil
	}

	// Request failed
	errorResponse, err := extractError(data)
	if err != nil {
		return nil, err
	}
	return nil, errorResponse
}
