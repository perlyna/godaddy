package godaddy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RequestExecute request execute
func (c *Context) RequestExecute(request *http.Request, result interface{}) error {
	if len(c.XShopperID) != 0 {
		request.Header.Set("X-Shopper-Id", c.XShopperID)
	}
	request.Header.Set("Authorization", c.authorization)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == 204 || response.StatusCode == 201 {
		return nil
	}
	if result == nil && response.StatusCode == 200 {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("response code : %d\nresponse body : %s", response.StatusCode, body)
	}
	if err = json.Unmarshal(body, result); err != nil {
		return err
	}
	return nil
}
