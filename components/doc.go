// Package components contains KubeVela ComponentDefinition implementations.
// Each component defines a workload type that can be used in Applications.
//
// To create a new component:
//
//	func init() {
//	    defkit.Register(MyComponent())
//	}
//
//	func MyComponent() *defkit.ComponentDefinition {
//	    return defkit.NewComponent("my-component").
//	        Description("My component description").
//	        Workload("apps/v1", "Deployment").
//	        // ... configuration
//	}
package components
