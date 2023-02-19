# Defining `discord` contact point types

## Example usage

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: AlertManager
metadata:
  name: alertmanager-example
spec:
  contact_points:
    - name: Team A

      # Contact Team A via Discord
      contacts:
        - discord:
            webhook: { value: "https://discord.com/api/webhooks/some_id/some_token" }

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
        - discord:
            # Webhook to use.
            # Required.
            webhook:
              # Webhook URL, as plain text. This is not recommended.
              # Optional. Default: ''
              value: ''

              # Reference to a secret containing the webhook URL.
              # Optional. Default: none
              valueFrom:
                secretKeyRef:
                  name: 'secret-name' # name of the secret
                  key: 'key' # Key within the secret

            # Use the username configured in Discord's webhook settings.
            # Otherwise, the username will be 'Grafana'.
            # Optional. Default: false
            use_discord_username: false
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)
