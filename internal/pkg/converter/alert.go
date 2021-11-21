package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertAlert(panel sdk.Panel) *grabana.Alert {
	if panel.Alert == nil {
		return nil
	}

	sdkAlert := panel.Alert

	notifications := make([]string, 0, len(sdkAlert.Notifications))
	for _, notification := range sdkAlert.Notifications {
		notifications = append(notifications, notification.UID)
	}

	alert := &grabana.Alert{
		Title:            sdkAlert.Name,
		Message:          sdkAlert.Message,
		EvaluateEvery:    sdkAlert.Frequency,
		For:              sdkAlert.For,
		Tags:             sdkAlert.AlertRuleTags,
		OnNoData:         sdkAlert.NoDataState,
		OnExecutionError: sdkAlert.ExecutionErrorState,
		Notifications:    notifications,
		If:               converter.convertAlertConditions(sdkAlert),
	}

	return alert
}

func (converter *JSON) convertAlertConditions(sdkAlert *sdk.Alert) []grabana.AlertCondition {
	conditions := make([]grabana.AlertCondition, 0, len(sdkAlert.Conditions))

	for _, condition := range sdkAlert.Conditions {
		conditions = append(conditions, grabana.AlertCondition{
			Operand: condition.Operator.Type,
			Value: grabana.AlertValue{
				Func:     condition.Reducer.Type,
				QueryRef: condition.Query.Params[0],
				From:     condition.Query.Params[1],
				To:       condition.Query.Params[2],
			},
			Threshold: converter.convertAlertThreshold(condition),
		})
	}

	return conditions
}

func (converter *JSON) convertAlertThreshold(sdkCondition sdk.AlertCondition) grabana.AlertThreshold {
	threshold := grabana.AlertThreshold{}

	switch sdkCondition.Evaluator.Type {
	case "no_value":
		threshold.HasNoValue = true
	case "lt":
		threshold.Below = &sdkCondition.Evaluator.Params[0]
	case "gt":
		threshold.Above = &sdkCondition.Evaluator.Params[0]
	case "outside_range":
		threshold.OutsideRange = [2]float64{sdkCondition.Evaluator.Params[0], sdkCondition.Evaluator.Params[1]}
	case "within_range":
		threshold.WithinRange = [2]float64{sdkCondition.Evaluator.Params[0], sdkCondition.Evaluator.Params[1]}
	}

	return threshold
}
