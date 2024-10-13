package utils

import (
	"context"
	"reflect"

	"errors"
	"fmt"
	"strconv"
	"strings"

	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/phamdinhha/event-booking-service/pkg/http_utils"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

func ParseUuidQuery(attribute string, uuid_str string) (uuid.UUID, *http_utils.AttributeError) {
	if uuid_str != "" {
		id, err := uuid.Parse(uuid_str)
		if err != nil {
			return uuid.Nil, &http_utils.AttributeError{
				Attribute: attribute,
				Cause:     err.Error(),
				Constraint: fmt.Sprintf(
					"%s must be a valid UUID v4 of length 36.",
					attribute,
				),
			}
		}
		return id, nil
	}
	return uuid.Nil, nil
}

func ValidateTimestampQuery(timestamp_str string) (int, error) {
	if timestamp_str != "" {
		timestamp, err := strconv.ParseInt(timestamp_str, 10, 64)
		if err != nil {
			return 0, err
		}
		if timestamp < 0 {
			return 0, errors.New("negative value")
		}
		return int(timestamp), nil
	}
	return 0, nil
}

func ParseTimestampQuery(
	query_map map[string]string,
) (map[string]int64, *http_utils.AttributeError) {
	out_map := map[string]int64{}

	for att, value := range query_map {
		ts, err := ValidateTimestampQuery(value)
		if err != nil {
			return nil, &http_utils.AttributeError{
				Attribute:  att,
				Cause:      err.Error(),
				Constraint: "timestamp must be a positive integer.",
			}
		}
		out_map[att] = int64(ts)
	}
	return out_map, nil
}

func ParseChoiceQuery(
	attribute string,
	choice string,
	accepted_choices []string,
) (string, *http_utils.AttributeError) {
	if choice != "" && !slices.Contains(accepted_choices, choice) {
		return "", &http_utils.AttributeError{
			Attribute:  attribute,
			Cause:      fmt.Sprintf("invalid choice: %s", choice),
			Constraint: fmt.Sprintf("choice must be one of: %s", strings.Join(accepted_choices, ", ")),
		}
	}
	return choice, nil
}

func ParseBool(
	attribute string,
	bool_str string,
) (bool, *http_utils.AttributeError) {

	if bool_str != "" {
		b, err := strconv.ParseBool(bool_str)
		if err != nil {
			return false, &http_utils.AttributeError{
				Attribute: attribute,
				Cause:     err.Error(),
				Constraint: fmt.Sprintf(
					"%s must be a valid boolean value.",
					attribute,
				),
			}
		}
		return b, nil
	}
	return false, nil
}

func ParseDate(
	attribute string,
	date_str string,
) (string, *http_utils.AttributeError) {
	//Initial format: "%d/%m/%Y"
	date_arr := strings.Split(date_str, "/")
	year_str := date_arr[2]
	month_str := date_arr[1]
	day_str := date_arr[0]

	year, err := strconv.Atoi(year_str)
	if err != nil {
		return "", &http_utils.AttributeError{
			Attribute: attribute,
			Cause:     err.Error(),
			Constraint: fmt.Sprintf(
				"%s must be a valid date of the format DD/MM/YYYY.",
				attribute,
			),
		}
	}

	month, err := strconv.Atoi(month_str)
	if err != nil {
		return "", &http_utils.AttributeError{
			Attribute: attribute,
			Cause:     err.Error(),
			Constraint: fmt.Sprintf(
				"%s must be a valid date of the format DD/MM/YYYY.",
				attribute,
			),
		}
	}

	day, err := strconv.Atoi(day_str)
	if err != nil {
		return "", &http_utils.AttributeError{
			Attribute: attribute,
			Cause:     err.Error(),
			Constraint: fmt.Sprintf(
				"%s must be a valid date of the format DD/MM/YYYY.",
				attribute,
			),
		}
	}

	time_to_return := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	return time_to_return, nil
}

func ValidateRequiredAttributes(
	T interface{},
) *http_utils.AttributeError {
	typ := reflect.TypeOf(T)
	val := reflect.ValueOf(T)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		field_tags := field.Tag
		if field_tags.Get("required") == "true" {
			if val.Field(i) == reflect.ValueOf(nil) {
				return &http_utils.AttributeError{
					Attribute:  field.Tag.Get("json"),
					Cause:      fmt.Sprintf("Missing field %s", field_tags.Get("json")),
					Constraint: fmt.Sprintf("Field %s is must be %s", field_tags.Get("json"), val.Type().Field(i).Name),
				}
			}
		}
	}

	return nil
}
