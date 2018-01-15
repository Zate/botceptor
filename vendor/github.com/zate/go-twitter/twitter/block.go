package twitter

import (
	"net/http"

	"github.com/dghubble/sling"
)

// BlockService provides methods for accessing Twitter user API endpoints for blocking other users.
type BlockService struct {
	sling *sling.Sling
}

// NewBlockService returns a new BlockService.
func newBlockService(sling *sling.Sling) *BlockService {
	return &BlockService{
		sling: sling.Path("blocks/"),
	}
}

// BlockUserParams is a struct for the params supplied to both Create and Destroy
type BlockUserParams struct {
	UserID          int64  `url:"user_id,omitempty"`          //optional - The screen name of the potentially blocked user. Helpful for disambiguating when a valid screen name is also a user ID.
	ScreenName      string `url:"screen_name,omitempty"`      // optional - The ID of the potentially blocked user. Helpful for disambiguating when a valid user ID is also a valid screen name.
	IncludeEntities *bool  `url:"include_entities,omitempty"` // optional - The entities node will not be included when set to false.  Default: false
	SkipStatus      *bool  `url:"skip_status,omitempty"`      // optional - When set to either true , t or 1 statuses will not be included in the returned user objects.  Default: true
}

// Create blocks a user
func (s *BlockService) Create(params *BlockUserParams) (*User, *http.Response, error) {
	user := new(User)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("create.json").QueryStruct(params).Receive(user, apiError)
	return user, resp, relevantError(err, *apiError)
}

// Destroy removes a block on a user
func (s *BlockService) Destroy(params *BlockUserParams) (*User, *http.Response, error) {
	user := new(User)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("destroy.json").QueryStruct(params).Receive(user, apiError)
	return user, resp, relevantError(err, *apiError)
}
