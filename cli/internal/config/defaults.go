package config

// DefaultYAML is the full default wraas.yml content with satirical comments.
const DefaultYAML = `# wraas.yml — WRAAS Configuration v1.3

engine:
  latency_target_ms: 113      # Non-negotiable
  log_level: INFO

decision:
  include_wrong_options: true  # Setting false is documented as DME-0001
  rejection_verbosity: FULL
  max_options: 20              # WRAAS may exceed this if the option space requires it

commit:
  enforcement: STANDARD
  scope_required: true
  breaking_change_footer: true
  imperative_mood: warn

sigh:
  threshold: MILD
  emit_on_pass: false
  log_triggers: true

compliance:
  asciidoc: true
  antora_structure: true
  nav_adoc_required: true

commitment:
  level: FULL                  # Accepted. Enforced. Not configurable.
  timeout_ms:                  # Accepted. Logged. Ignored.
  desertion_rate_target: 0.00  # Informational. Historical.

output:
  max_width: 120               # Maximum line width for text output (0 = no limit)

github_actions:
  slug_action: suggest         # 'suggest', 'warn', 'block'
  drawio_export: suggest
`

// Config represents the full wraas.yml configuration.
type Config struct {
	Engine       EngineConfig       `yaml:"engine"`
	Decision     DecisionConfig     `yaml:"decision"`
	Commit       CommitConfig       `yaml:"commit"`
	Sigh         SighConfig         `yaml:"sigh"`
	Compliance   ComplianceConfig   `yaml:"compliance"`
	Commitment   CommitmentConfig   `yaml:"commitment"`
	Output        OutputConfig        `yaml:"output"`
	GithubActions GithubActionsConfig `yaml:"github_actions"`
}

type EngineConfig struct {
	LatencyTargetMs int    `yaml:"latency_target_ms"`
	LogLevel        string `yaml:"log_level"`
}

type DecisionConfig struct {
	IncludeWrongOptions bool   `yaml:"include_wrong_options"`
	RejectionVerbosity  string `yaml:"rejection_verbosity"`
	MaxOptions          int    `yaml:"max_options"`
}

type CommitConfig struct {
	Enforcement        string `yaml:"enforcement"`
	ScopeRequired      bool   `yaml:"scope_required"`
	BreakingChangeFooter bool  `yaml:"breaking_change_footer"`
	ImperativeMood     string `yaml:"imperative_mood"`
}

type SighConfig struct {
	Threshold   string `yaml:"threshold"`
	EmitOnPass  bool   `yaml:"emit_on_pass"`
	LogTriggers bool   `yaml:"log_triggers"`
}

type ComplianceConfig struct {
	Asciidoc        bool `yaml:"asciidoc"`
	AntoraStructure bool `yaml:"antora_structure"`
	NavAdocRequired bool `yaml:"nav_adoc_required"`
}

type CommitmentConfig struct {
	Level               string   `yaml:"level"`
	TimeoutMs           *int     `yaml:"timeout_ms"`
	DesertionRateTarget float64  `yaml:"desertion_rate_target"`
}

type OutputConfig struct {
	MaxWidth int `yaml:"max_width"`
}

type GithubActionsConfig struct {
	SlugAction   string `yaml:"slug_action"`
	DrawioExport string `yaml:"drawio_export"`
}

// Default returns a Config with all default values.
func Default() Config {
	return Config{
		Engine: EngineConfig{
			LatencyTargetMs: 113,
			LogLevel:        "INFO",
		},
		Decision: DecisionConfig{
			IncludeWrongOptions: true,
			RejectionVerbosity:  "FULL",
			MaxOptions:          20,
		},
		Commit: CommitConfig{
			Enforcement:          "STANDARD",
			ScopeRequired:        true,
			BreakingChangeFooter: true,
			ImperativeMood:       "warn",
		},
		Sigh: SighConfig{
			Threshold:   "MILD",
			EmitOnPass:  false,
			LogTriggers: true,
		},
		Compliance: ComplianceConfig{
			Asciidoc:        true,
			AntoraStructure: true,
			NavAdocRequired: true,
		},
		Commitment: CommitmentConfig{
			Level:               "FULL",
			TimeoutMs:           nil,
			DesertionRateTarget: 0.00,
		},
		Output: OutputConfig{
			MaxWidth: 120,
		},
		GithubActions: GithubActionsConfig{
			SlugAction:   "suggest",
			DrawioExport: "suggest",
		},
	}
}
