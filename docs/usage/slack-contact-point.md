# Defining `slack` contact point types

## Example usage

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: AlertManager
metadata:
  name: alertmanager-example
spec:
  contact_points:
    - name: Team A
      
      # Contact Team A via email
      contacts:
        - slack:
            webhook:
              value: 'https://api.slack.com/messaging/webhooks'

  routing:
    # ... omitted
```

## Reference

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: AlertManager
metadata:
  name: alertmanager-example
spec:
  # Alerts not matched by any of the routing rules will be sent to this contact point
  default_contact_point: 'Team A'

  # List of known contact points
  contact_points:
    - name: Team A
      contacts:
        - slack:
            # Slack webhook Url.
            # Required.
            webhook:
              # Webhook, as plain text. This is not recommended.
              # Optional. Default: ''
              value: ''
        
              # Reference to a secret containing the webhook URL.
              # Optional. Default: none
              valueFrom:
                secretKeyRef:
                  name: 'secret-name' # name of the secret
                  key: 'url' # Key within the secret

            # Templated title of the slack message.
            # Optional. Default: ''
            title: ''
            
            # Body of the slack message.
            # Optional. Default: ''
            body: ''
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)