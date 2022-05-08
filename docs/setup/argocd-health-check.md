# Integrating with ArgoCD health checks

ArgoCD supports [health checks for custom resources](https://argo-cd.readthedocs.io/en/stable/operator-manual/health/#way-1-define-a-custom-health-check-in-argocd-cm-configmap).
To enable it for Dark-managed manifests, add the following code to your `argo-cm` ConfigMap:

```yaml
data:
  resource.customizations.health.k8s.kevingomez.fr_GrafanaDashboard: |
    hs = {}                                                                                                                                                                                                          
    if obj.status ~= nil then                                                                                                                                                                                        
      if obj.status.status ~= "OK" then                                                                                                                                                                              
        hs.status = "Degraded"                                                                                                                                                                                       
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      else                                                                                                                                                                                                           
        hs.status = "Healthy"                                                                                                                                                                                        
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      end                                                                                                                                                                                                            
    end                                                                                                                                                                                                              

    hs.status = "Progressing"                                                                                                                                                                                            
    hs.message = "Status unknown"                                                                                                                                                                                    
    return hs
  resource.customizations.health.k8s.kevingomez.fr_Datasource: |
    hs = {}                                                                                                                                                                                                          
    if obj.status ~= nil then                                                                                                                                                                                        
      if obj.status.status ~= "OK" then                                                                                                                                                                              
        hs.status = "Degraded"                                                                                                                                                                                       
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      else                                                                                                                                                                                                           
        hs.status = "Healthy"                                                                                                                                                                                        
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      end                                                                                                                                                                                                            
    end                                                                                                                                                                                                              

    hs.status = "Progressing"                                                                                                                                                                                            
    hs.message = "Status unknown"                                                                                                                                                                                    
    return hs
  resource.customizations.health.k8s.kevingomez.fr_APIKey: |
    hs = {}                                                                                                                                                                                                          
    if obj.status ~= nil then                                                                                                                                                                                        
      if obj.status.status ~= "OK" then                                                                                                                                                                              
        hs.status = "Degraded"                                                                                                                                                                                       
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      else                                                                                                                                                                                                           
        hs.status = "Healthy"                                                                                                                                                                                        
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      end                                                                                                                                                                                                            
    end                                                                                                                                                                                                              

    hs.status = "Progressing"                                                                                                                                                                                            
    hs.message = "Status unknown"                                                                                                                                                                                    
    return hs
  resource.customizations.health.k8s.kevingomez.fr_AlertManager: |
    hs = {}                                                                                                                                                                                                          
    if obj.status ~= nil then                                                                                                                                                                                        
      if obj.status.status ~= "OK" then                                                                                                                                                                              
        hs.status = "Degraded"                                                                                                                                                                                       
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      else                                                                                                                                                                                                           
        hs.status = "Healthy"                                                                                                                                                                                        
        hs.message = obj.status.message                                                                                                                                                                              
        return hs                                                                                                                                                                                                    
      end                                                                                                                                                                                                            
    end                                                                                                                                                                                                              

    hs.status = "Progressing"                                                                                                                                                                                            
    hs.message = "Status unknown"                                                                                                                                                                                    
    return hs
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)