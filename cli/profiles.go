package main

// Profiles describe an intended approach; selecting one never scaffolds a project.
type architectureProfile struct {
	ID, Name, Description, NameES, DescriptionES string
	Hybrid                                       bool
}

func localized(es bool, en, spanish string) string {
	if es {
		return spanish
	}
	return en
}
func architectureProfiles() []architectureProfile {
	return []architectureProfile{
		{"minimal", "Minimal monolith", "One application with simple boundaries; a low-overhead starting point.", "Monolito mínimo", "Una aplicación con límites sencillos; punto de partida de bajo costo.", false},
		{"layered", "Layered / N-tier", "Separate presentation, application and data responsibilities into layers.", "Capas / N niveles", "Separa presentación, aplicación y datos en capas de responsabilidad.", false},
		{"vertical-slices", "Vertical slices", "Organize each use case end to end, rather than by technical layer.", "Cortes verticales", "Organiza cada caso de uso de extremo a extremo, no por capa técnica.", false},
		{"clean", "Clean Architecture", "Keep domain and use cases independent of frameworks and infrastructure.", "Arquitectura limpia", "Mantiene dominio y casos de uso independientes de frameworks e infraestructura.", false},
		{"hexagonal", "Hexagonal / ports and adapters", "Isolate application behavior behind ports with replaceable external adapters.", "Hexagonal / puertos y adaptadores", "Aísla el comportamiento mediante puertos y adaptadores externos intercambiables.", false},
		{"onion", "Onion Architecture", "Place domain rules at the center, with dependencies directed inward.", "Arquitectura de cebolla", "Sitúa las reglas de dominio en el centro y dirige las dependencias hacia dentro.", false},
		{"modular-monolith", "Modular monolith", "One deployment with explicit business-module boundaries and ownership.", "Monolito modular", "Un despliegue con límites explícitos y responsabilidades por módulo de negocio.", false},
		{"microservices", "Microservices", "Independently deployed services; adds distributed operations and integration costs.", "Microservicios", "Servicios desplegables de forma independiente; añade costos de operación e integración distribuida.", false},
		{"event-driven", "Event-driven", "Components react to published events; requires clear delivery and consistency rules.", "Orientada a eventos", "Los componentes reaccionan a eventos; requiere reglas claras de entrega y consistencia.", false},
		{"serverless", "Serverless functions", "Managed functions for workloads; a deployment style, not a domain design by itself.", "Funciones sin servidor", "Funciones administradas para cargas de trabajo; es un estilo de despliegue, no un diseño de dominio por sí solo.", false},
		{"mvc", "Model-View-Controller", "Separate state, rendering and request/interaction coordination.", "Modelo-Vista-Controlador", "Separa estado, presentación y coordinación de solicitudes o interacciones.", false},
		{"mvvm", "Model-View-ViewModel", "Expose presentation state and behavior through view models for the UI.", "Modelo-Vista-Modelo de vista", "Expone el estado y comportamiento de presentación mediante modelos de vista.", false},
		{"component-ui", "Component-based UI", "Compose reusable UI components with explicit state and interaction boundaries.", "Interfaz basada en componentes", "Compone componentes reutilizables con límites explícitos de estado e interacción.", false},
		{"feature-ui", "Feature-based UI", "Group UI, state and API integration by product feature.", "Interfaz organizada por funcionalidades", "Agrupa interfaz, estado e integración con API por funcionalidad del producto.", false},
		{"modular-vertical-slices", "Modular monolith + vertical slices", "Keep module boundaries while organizing each module around complete use cases.", "Monolito modular + cortes verticales", "Mantiene límites entre módulos y organiza cada uno alrededor de casos de uso completos.", true},
		{"clean-vertical-slices", "Clean + vertical slices", "Use inward dependency rules while grouping application work by use case.", "Arquitectura limpia + cortes verticales", "Combina dependencias hacia el dominio con organización por caso de uso.", true},
		{"hexagonal-ddd", "Hexagonal + domain-driven design", "Combine ports/adapters with explicit domain models when business complexity warrants them.", "Hexagonal + diseño guiado por el dominio", "Combina puertos y adaptadores con modelos de dominio cuando la complejidad de negocio lo justifica.", true},
		{"microservices-events", "Microservices + events", "Use events between independently deployed services; plan failure and consistency handling.", "Microservicios + eventos", "Usa eventos entre servicios independientes; requiere planificar fallos y consistencia.", true},
		{"cqrs-event-sourcing", "CQRS + event sourcing", "Separate read/write models and retain an event history; a pattern combination with added complexity.", "CQRS + registro de eventos", "Separa modelos de lectura y escritura y conserva un historial de eventos; combinación de patrones con mayor complejidad.", true},
		{"custom", "Custom architecture", "Describe another approach or hybrid in your own words.", "Arquitectura personalizada", "Describe otro enfoque o una combinación con tus propias palabras.", false},
	}
}
func architectureOptions(es bool) []pickerOption {
	options := []pickerOption{}
	for _, p := range architectureProfiles() {
		kind := localized(es, "Architecture", "Arquitectura")
		if p.Hybrid {
			kind = localized(es, "Hybrid", "Híbrida")
		}
		options = append(options, pickerOption{ID: p.ID, Label: localized(es, p.Name, p.NameES), Description: localized(es, p.Description, p.DescriptionES), Topics: []string{kind}})
	}
	return options
}
func scopeOptions(es bool) []pickerOption {
	return []pickerOption{
		{ID: "frontend", Label: "Frontend", Description: localized(es, "Web interfaces, browser clients and presentation behavior.", "Interfaces web, clientes de navegador y comportamiento de presentación.")},
		{ID: "backend", Label: "Backend", Description: localized(es, "APIs, server-side applications and business/data operations.", "API, aplicaciones del servidor y operaciones de negocio y datos.")},
		{ID: "app", Label: localized(es, "App (mobile / hybrid)", "App (móvil / híbrida)"), Description: localized(es, "Mobile or hybrid applications and native-device integration.", "Aplicaciones móviles o híbridas e integración con capacidades del dispositivo.")},
	}
}
