package hermod

import "encoding/json"

func (r *Response) JSON(v any) error {
	return json.Unmarshal(r.Data, v)
}
