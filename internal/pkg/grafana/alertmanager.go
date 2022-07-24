package grafana

import (
	"context"
	"fmt"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/grabana/alertmanager/discord"
	"github.com/K-Phoen/grabana/alertmanager/email"
	"github.com/K-Phoen/grabana/alertmanager/opsgenie"
	"github.com/K-Phoen/grabana/alertmanager/slack"
	"github.com/go-logr/logr"
)

var ErrInvalidContactPointType = fmt.Errorf("invalid contact point type")
var ErrInvalidRoutingRule = fmt.Errorf("invalid routing rule")

type AlertManager struct {
	logger        logr.Logger
	grabanaClient *grabana.Client
	refReader     refReader
}

func NewAlertManager(logger logr.Logger, grabanaClient *grabana.Client, refReader refReader) *AlertManager {
	return &AlertManager{
		logger:        logger,
		grabanaClient: grabanaClient,
		refReader:     refReader,
	}
}

func (manager *AlertManager) Reset(ctx context.Context) error {
	config := alertmanager.New(
		alertmanager.ContactPoints(
			alertmanager.ContactPoint("grafana-default-email", email.To([]string{"<example@email.com>"})),
		),
		alertmanager.DefaultContactPoint("grafana-default-email"),
	)

	return manager.grabanaClient.ConfigureAlertManager(ctx, config)
}

func (manager *AlertManager) Configure(ctx context.Context, manifest v1alpha1.AlertManager) error {
	var managerOpts []alertmanager.Option

	// message templates
	if len(manifest.Spec.MessageTemplates) != 0 {
		managerOpts = append(managerOpts, alertmanager.Templates(manifest.Spec.MessageTemplates))
	}

	// contact points
	contactPointsOpts, err := manager.contactPointsOpts(ctx, manifest)
	if err != nil {
		return err
	}
	if len(contactPointsOpts) != 0 {
		managerOpts = append(managerOpts, alertmanager.ContactPoints(contactPointsOpts...))
	}

	// routing rules
	routingOpts, err := manager.routingOpts(manifest)
	if err != nil {
		return err
	}
	if len(routingOpts) != 0 {
		managerOpts = append(managerOpts, alertmanager.Routing(routingOpts...))
	}

	// default contact point
	if manifest.Spec.DefaultContactPoint != "" {
		managerOpts = append(managerOpts, alertmanager.DefaultContactPoint(manifest.Spec.DefaultContactPoint))
	}

	// default grouping labels
	if len(manifest.Spec.DefaultGroupBy) != 0 {
		managerOpts = append(managerOpts, alertmanager.DefaultGroupBys(manifest.Spec.DefaultGroupBy...))
	}

	return manager.grabanaClient.ConfigureAlertManager(ctx, alertmanager.New(managerOpts...))
}

func (manager *AlertManager) contactPointsOpts(ctx context.Context, manifest v1alpha1.AlertManager) ([]alertmanager.Contact, error) {
	opts := []alertmanager.Contact{}

	for _, contactPointSpec := range manifest.Spec.ContactPoints {
		opt, err := manager.contactPointOpt(ctx, manifest.Namespace, contactPointSpec)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (manager *AlertManager) contactPointOpt(ctx context.Context, namespace string, contactPointSpec v1alpha1.ContactPoint) (alertmanager.Contact, error) {
	contactPointTypeOpts := []alertmanager.ContactPointOption{}

	for _, contactPointType := range contactPointSpec.Contacts {
		opt, err := manager.contactPointTypeOpt(ctx, namespace, contactPointType)
		if err != nil {
			return alertmanager.Contact{}, err
		}

		contactPointTypeOpts = append(contactPointTypeOpts, opt)
	}

	return alertmanager.ContactPoint(contactPointSpec.Name, contactPointTypeOpts...), nil
}

func (manager *AlertManager) contactPointTypeOpt(ctx context.Context, namespace string, contactPointType v1alpha1.ContactPointType) (alertmanager.ContactPointOption, error) {
	if contactPointType.Email != nil {
		return manager.contactPointTypeEmail(*contactPointType.Email)
	}
	if contactPointType.Slack != nil {
		return manager.contactPointTypeSlack(ctx, namespace, *contactPointType.Slack)
	}
	if contactPointType.Opsgenie != nil {
		return manager.contactPointTypeOpsgenie(ctx, namespace, *contactPointType.Opsgenie)
	}
	if contactPointType.Discord != nil {
		return manager.contactPointTypeDiscord(ctx, namespace, *contactPointType.Discord)
	}

	return nil, ErrInvalidContactPointType
}

func (manager *AlertManager) contactPointTypeEmail(contactPointType v1alpha1.EmailContactType) (alertmanager.ContactPointOption, error) {
	opts := []email.Option{}

	if contactPointType.Single {
		opts = append(opts, email.Single())
	}
	if contactPointType.Message != "" {
		opts = append(opts, email.Message(contactPointType.Message))
	}

	return email.To(contactPointType.To, opts...), nil
}

func (manager *AlertManager) contactPointTypeSlack(ctx context.Context, namespace string, contactPointType v1alpha1.SlackContactType) (alertmanager.ContactPointOption, error) {
	opts := []slack.Option{}

	webhookURL, err := manager.refReader.RefToValue(ctx, namespace, contactPointType.Webhook)
	if err != nil {
		return nil, err
	}

	if contactPointType.Title != "" {
		opts = append(opts, slack.Title(contactPointType.Title))
	}
	if contactPointType.Body != "" {
		opts = append(opts, slack.Body(contactPointType.Body))
	}

	return slack.Webhook(webhookURL, opts...), nil
}

func (manager *AlertManager) contactPointTypeOpsgenie(ctx context.Context, namespace string, contactPointType v1alpha1.OpsgenieContactType) (alertmanager.ContactPointOption, error) {
	opts := []opsgenie.Option{}

	apiKey, err := manager.refReader.RefToValue(ctx, namespace, contactPointType.APIKey)
	if err != nil {
		return nil, err
	}

	if contactPointType.OverridePriority {
		opts = append(opts, opsgenie.OverridePriority())
	}
	if contactPointType.AutoClose {
		opts = append(opts, opsgenie.AutoClose())
	}

	return opsgenie.With(contactPointType.APIURL, apiKey, opts...), nil
}

func (manager *AlertManager) contactPointTypeDiscord(ctx context.Context, namespace string, contactPointType v1alpha1.DiscordContactType) (alertmanager.ContactPointOption, error) {
	opts := []discord.Option{}

	webhook, err := manager.refReader.RefToValue(ctx, namespace, contactPointType.Webhook)
	if err != nil {
		return nil, err
	}

	if contactPointType.UseDiscordUsername {
		opts = append(opts, discord.UseDiscordUsername())
	}

	return discord.With(webhook, opts...), nil
}

func (manager *AlertManager) routingOpts(manifest v1alpha1.AlertManager) ([]alertmanager.RoutingPolicy, error) {
	opts := []alertmanager.RoutingPolicy{}

	for _, routingPolicy := range manifest.Spec.Routing {
		opt, err := manager.routingPolicyOpt(routingPolicy)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (manager *AlertManager) routingPolicyOpt(policySpec v1alpha1.RoutingPolicy) (alertmanager.RoutingPolicy, error) {
	opts := []alertmanager.RoutingPolicyOption{}

	for _, rule := range policySpec.Rules {
		labelOpts, err := manager.routingLabelRules(rule)
		if err != nil {
			return alertmanager.RoutingPolicy{}, err
		}

		opts = append(opts, labelOpts...)
	}

	return alertmanager.Policy(policySpec.ContactPoint, opts...), nil
}

func (manager *AlertManager) routingLabelRules(rule v1alpha1.LabelsMatchingRule) ([]alertmanager.RoutingPolicyOption, error) {
	var operatorFunc func(tag string, value string) alertmanager.RoutingPolicyOption
	var labels map[string]string
	opts := []alertmanager.RoutingPolicyOption{}

	//nolint:gocritic
	if rule.Eq != nil {
		operatorFunc = alertmanager.TagEq
		labels = rule.Eq
	} else if rule.Neq != nil {
		operatorFunc = alertmanager.TagNeq
		labels = rule.Neq
	} else if rule.Matches != nil {
		operatorFunc = alertmanager.TagMatches
		labels = rule.Matches
	} else if rule.NotMatches != nil {
		operatorFunc = alertmanager.TagNotMatches
		labels = rule.NotMatches
	} else {
		return nil, ErrInvalidRoutingRule
	}

	for key, value := range labels {
		opts = append(opts, operatorFunc(key, value))
	}

	return opts, nil
}
