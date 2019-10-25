package subscription

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	subscription struct {
		sync.RWMutex

		domain        string
		expires       time.Time
		limitMaxUsers uint
		isTrial       bool
		isValid       bool
	}

	PermitChecker interface {
		Validate(string, bool) error
		CanCreateUser(uint) error
		CanRegister(uint) error
	}
)

var (
	now = func() time.Time {
		return time.Now()
	}
)

const (
	salesEmail = "sales@crust.tech"

	// If less then this amount of days we'll show admins a warning
	warnAdminDaysLimit = 30

	// If less then this amount of days we'll show everyone using a trial a warning
	warnTrialDaysLimit = 10

	// If we got trial subscription w/o max-user limit, this
	// is the fallback number
	limitMaxUsersTrialDefault = 10

	// See subscription struct's functions on how & where these strings are used
	willExpire      = `This Crust subscription will expire on [exp-date]. Please contact [sales-email] to renew the subscription.`
	hasExpired      = `This Crust subscription has expired. Please contact your administrator or [sales-email] to renew the subscription.`
	trialWillExpire = `This Crust trial will expire on [exp-date]. To convert this trial in a to a subscription, please contact [sales-email].`
	trialHasExpired = `Your Crust trial has expired. Please contact [sales-email] to learn how to convert this trial in to a Crust subscription.`
	invalidKey      = `Unverified or invalid subscription key. Please contact your administrator or [sales-email].`

	trialAddUserError = `The Crust trial is limited to [user-limit] user(s). If you need more users, please contact us at [sales-email].`
	addUserError      = `Your subscription user limit has been reached. Please contact [sales-email] to learn how to increase the number of users.`
	signupError       = `Registration is disabled at the moment. Please contact your administrator.`
)

func New(domain string) *subscription {
	return &subscription{domain: domain}
}

// Update updates subscription data with new values
//
// We're keeping only values revevant to Validate & Can*() functions
func (s *subscription) Update(expires time.Time, isTrial bool, limitMaxUsers uint) {
	s.Lock()
	defer s.Unlock()

	s.expires = expires
	s.isTrial = isTrial
	s.isValid = s.domain != "" && !s.expires.IsZero()

	if limitMaxUsers == 0 && isTrial {
		// Trial w/o user limit?
		// set to default
		s.limitMaxUsers = limitMaxUsersTrialDefault
	} else {
		s.limitMaxUsers = limitMaxUsers
	}
}

func (s *subscription) Reset() {
	s.Lock()
	defer s.Unlock()

	s.expires = time.Time{}
	s.limitMaxUsers = 0
	s.isTrial = false
	s.isValid = false
}

// Validate checks domain and expiration date
//
// It returns different kinds of errors, depending on current state of
// subscription:
//   - is trial
//   - number of days from expiration
//   - isAdmin flag
//
// All states that this function validates are technically not errors
// but just warnings. Nevertheless we always return an error for consistency
// and simpler func signature
func (s *subscription) Validate(domain string, isAdmin bool) error {
	s.RLock()
	defer s.RUnlock()

	var (
		// How many days are left until this subscription expires?
		daysLeft = math.Floor(s.expires.Sub(now()).Hours() / 24)
	)

	switch true {
	case !s.isValid:
		return s.error(invalidKey)

	case s.isTrial && daysLeft <= warnTrialDaysLimit:
		return s.error(trialWillExpire)

	case s.isTrial && daysLeft <= 0:
		return s.error(trialHasExpired)

	case daysLeft <= warnAdminDaysLimit && isAdmin:
		return s.error(willExpire)

	case daysLeft <= 0:
		return s.error(hasExpired)

	default:
		return nil
	}
}

// CanCreateUser - Does subscription allow us to create new user
func (s *subscription) CanCreateUser(currentTotal uint) error {
	s.Lock()
	defer s.Unlock()

	var (
		// Compare limit with current total if user limit is set (> 0)
		overLimit = s.limitMaxUsers > 0 && currentTotal >= s.limitMaxUsers
	)

	switch true {
	case !s.isValid:
		return s.error(invalidKey)

	case s.isTrial && overLimit:
		return s.error(trialAddUserError)

	case overLimit:
		return s.error(addUserError)

	default:
		return nil
	}
}

// CanRegister - Can users (self) register (same rules as CanCreateUser but different error)
//
// We'll be showing this to everyone, so let's be careful not to tell too much
func (s *subscription) CanRegister(currentTotal uint) error {
	s.Lock()
	defer s.Unlock()

	if s.isValid && s.limitMaxUsers == 0 || currentTotal < s.limitMaxUsers {
		return nil
	}

	return s.error(signupError)
}

// Converts error template into error using permit values
func (s subscription) error(t string) error {
	t = strings.NewReplacer(
		"[exp-date]", s.expires.Format(time.RFC1123),
		"[sales-email]", salesEmail,
		"[user-limit]", strconv.Itoa(int(s.limitMaxUsers)),
	).Replace(t)

	return errors.New(t)
}

var _ PermitChecker = &subscription{}
