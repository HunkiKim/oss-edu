kubectl get pods -A -o=jsonpath="{range .items[*]}{.metadata.name}{\"\t\"}{.metadata.uid}{\"\t\"}{.metadata.namespace}{\"\n\"}{end}"
