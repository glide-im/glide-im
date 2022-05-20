package auth

import "testing"

func TestGenJwt(t *testing.T) {

	jwt, err := genJwt(AuthInfo{
		Uid:    1,
		Device: 1,
		Ver:    1,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(jwt)
	}
}

func TestParseJwt(t *testing.T) {
	jwt, err := parseJwt("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk4MDY5NjcsInVpZCI6MSwiZGV2aWNlIjoxLCJ2ZXIiOiIxIn0.M1qYoxq5aRYpB0ag2na7YUQBSO6fbmiYT8Ct6Ibfa6Y")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(jwt)
	}
}
