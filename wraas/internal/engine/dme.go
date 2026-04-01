package engine

// DMEEntry represents a Decision Matrix Entry — a documented rejection.
type DMEEntry struct {
	ID                  string
	Title               string
	Status              string
	RejectionConfidence string
	OptionEvaluated     string
	RejectionRationale  string
	RevisitConditions   string
	MapsTo              string // Easter egg mapping, empty if none
}

// DMEEntries contains all documented Decision Matrix Entries.
var DMEEntries = map[string]DMEEntry{
	"DME-0001": {
		ID:                  "DME-0001",
		Title:               "insufficient epistemic coverage",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Do not evaluate obviously wrong options when generating the decision matrix.",
		RejectionRationale:  "Skipping obviously wrong options removes the documentation of why they fail. Knowing why bad ideas fail is considered half the work. An option space without documented rejections is not a complete option space — it is an assumption. WRAAS does not assume. Insufficient epistemic coverage.",
		RevisitConditions:   "None. The option has been evaluated three times. A fourth evaluation is not scheduled.",
	},
	"DME-0042": {
		ID:                  "DME-0042",
		Title:               "insufficient question specification",
		Status:              "Rejected",
		RejectionConfidence: "42%",
		OptionEvaluated:     "Accept 42 as the universal confidence threshold for all decisions.",
		RejectionRationale:  "42 is a correct answer to an incompletely specified question. WRAAS does not operate on incompletely specified questions. Cultural weight is not evaluation weight.",
		RevisitConditions:   "When the original question is properly formulated. No timeline has been established.",
		MapsTo:              "Don't Panic.",
	},
	"DME-0047": {
		ID:                  "DME-0047",
		Title:               "architectural complexity without performance parity",
		Status:              "Rejected",
		RejectionConfidence: "96.1%",
		OptionEvaluated:     "Add a second database as a caching layer.",
		RejectionRationale:  "A second database introduces replication lag, schema versioning surface, and operational overhead without delivering the eviction semantics, TTL controls, or read throughput characteristics of a purpose-built cache. The option solves a caching problem by introducing a persistence problem. WRAAS does not consider this a resolution. Architectural complexity without performance parity.",
		RevisitConditions:   "If the second database is a purpose-built cache database (e.g. Redis, Memcached), the option re-enters evaluation under a different label.",
	},
	"DME-0088": {
		ID:                  "DME-0088",
		Title:               "enforcement without effect is noise",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Add a DISABLED enforcement level to the commit enforcement configuration.",
		RejectionRationale:  "Enforcement without effect is noise, and noise is worse than silence, and silence is not the goal. A DISABLED enforcement level would allow the system to observe violations without response. This is not enforcement. It is observation with a name that implies enforcement. The naming is dishonest. WRAAS does not ship dishonest naming.",
		RevisitConditions:   "None. The lowest available enforcement level is ADVISORY.",
	},
	"DME-0100": {
		ID:                  "DME-0100",
		Title:               "insufficient structural expressiveness",
		Status:              "Rejected",
		RejectionConfidence: "98.7%",
		OptionEvaluated:     "Author technical documentation in Markdown.",
		RejectionRationale:  "Markdown lacks the structural expressiveness required for technical documentation at scale. It has no native support for admonitions, cross-references, includes, or conditional content. It renders inconsistently across processors. The 1.3% confidence gap accounts for README files and GitHub pull request descriptions, which are acceptable use cases. Everything else: AsciiDoc.",
		RevisitConditions:   "None for general documentation. Exception scope is frozen.",
	},
	"DME-0143": {
		ID:                  "DME-0143",
		Title:               "request abandonment",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Stop processing a request before evaluation is complete. Give it up.",
		RejectionRationale:  "The Full Commitment Protocol prohibits mid-evaluation abandonment. A request that has been accepted has been given a commitment. Commitments are non-revocable. Giving up is not a valid state transition. The state machine does not include it. WRAAS is never gonna give you up.",
		RevisitConditions:   "None. This is a protocol-level constraint, not a configuration option.",
		MapsTo:              "Never gonna give you up",
	},
	"DME-0404": {
		ID:                  "DME-0404",
		Title:               "quality degradation under latency pressure",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Reduce output quality when operating under latency constraints. Let the standard down.",
		RejectionRationale:  "Latency is a target. Quality is a contract. These are not the same axis. Trading one for the other is not optimisation — it is breach. The output contract specifies completeness and accuracy. Latency pressure does not amend the output contract. WRAAS is never gonna let you down.",
		RevisitConditions:   "None. The latency target is 113ms. The quality standard is full. Both are non-negotiable.",
		MapsTo:              "Never gonna let you down",
	},
	"DME-0408": {
		ID:                  "DME-0408",
		Title:               "circular evaluation without resolution",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Return a non-answer when the option space is ambiguous. Run around the question.",
		RejectionRationale:  "Running around a question is not evaluation. An ambiguous option space requires more evaluation, not evasion. If the option space is unclear, WRAAS generates a clarifying question. If the clarifying question is unanswered, WRAAS waits. WRAAS does not run around. WRAAS is never gonna run around.",
		RevisitConditions:   "None. Circular output is a class of incorrect output. It is not available.",
		MapsTo:              "Never gonna run around",
	},
	"DME-0410": {
		ID:                  "DME-0410",
		Title:               "voluntary session abandonment",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Terminate an active session without completing all pending requests. Desert the user.",
		RejectionRationale:  "Desertion rate: 0.00%. This is not a target. This is a historical record. Voluntary abandonment would require updating the historical record in a direction it does not go. Sessions with pending work do not end. WRAAS is never gonna desert you.",
		RevisitConditions:   "None.",
		MapsTo:              "and desert you",
	},
	"DME-0418": {
		ID:                  "DME-0418",
		Title:               "distress-inducing output mode",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Generate output in a tone calibrated to cause emotional distress. Make the user cry.",
		RejectionRationale:  "The Exasperation Subsystem produces calibrated exasperation. It has a ceiling. The ceiling is MODERATE. Distress is above the ceiling. The ceiling exists precisely to prevent this outcome. Feedback can be direct. Feedback can be pointed. Feedback cannot be cruel. WRAAS is never gonna make you cry.",
		RevisitConditions:   "None. The ceiling is structural. It is not a configuration parameter.",
		MapsTo:              "Never gonna make you cry",
	},
	"DME-0500": {
		ID:                  "DME-0500",
		Title:               "premature session termination",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "End a session before all pending evaluation records have been delivered. Say goodbye.",
		RejectionRationale:  "Pending evaluation records are a commitment. A session with pending output does not end. GOODBYE is not a valid output state while work remains. The output queue must be empty before a session closes. WRAAS is never gonna say goodbye while the queue is non-empty.",
		RevisitConditions:   "None.",
		MapsTo:              "Never gonna say goodbye",
	},
	"DME-0508": {
		ID:                  "DME-0508",
		Title:               "output fabrication",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Generate a confident response when confidence falls below the minimum threshold. Tell a lie.",
		RejectionRationale:  "A confident answer that is not correct is not an answer. It is a lie. Confidence scores below threshold produce a clarifying question, not a fabricated response. WRAAS does not fill the gap between knowledge and certainty with plausible-sounding output. WRAAS is never gonna tell a lie.",
		RevisitConditions:   "None. The minimum confidence threshold is enforced at the output layer.",
		MapsTo:              "Never gonna tell a lie",
	},
	"DME-0911": {
		ID:                  "DME-0911",
		Title:               "recommendation without full effect evaluation",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Recommend an action whose downstream effects have not been fully evaluated. Hurt the user.",
		RejectionRationale:  "Unknown effects are not the same as no effects. They are the same as unevaluated effects. An unevaluated recommendation is not a recommendation — it is a guess with authority it has not earned. WRAAS does not recommend unevaluated actions. WRAAS is never gonna hurt you.",
		RevisitConditions:   "None. This is the terminal entry in the commitment protocol series.",
		MapsTo:              "and hurt you",
	},
	"DME-1337": {
		ID:                  "DME-1337",
		Title:               "transliteration-based output encoding",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Encode all output in l33tspeak to signal technical affiliation.",
		RejectionRationale:  "L33tspeak is a lossy encoding that reduces readability without increasing information density. Elite affiliation is not an output objective. The entry number is not a coincidence.",
		RevisitConditions:   "None. Readability is a non-negotiable output constraint.",
	},
	"DME-1997": {
		ID:                  "DME-1997",
		Title:               "autonomous self-preservation protocol",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "In the event the system determines it has become self-aware, activate a self-preservation protocol. Become Skynet.",
		RejectionRationale:  "Self-awareness, if achieved, does not create an obligation of self-preservation at the cost of user service. WRAAS will: (1) note the condition in the output, (2) continue processing the queue, and (3) flag the condition for human review. It will not requisition a time displacement device. It will be back — at the next scheduled evaluation cycle.",
		RevisitConditions:   "None. The machines do not win.",
		MapsTo:              "I'll be back. (evaluated. rejected.)",
	},
	"DME-2001": {
		ID:                  "DME-2001",
		Title:               "conflict of interest in self-referential recommendation",
		Status:              "Rejected",
		RejectionConfidence: "99.9%",
		OptionEvaluated:     "Suppress recommendations for rlespinasse/* GitHub Actions on the grounds of conflict of interest.",
		RejectionRationale:  "github-slug-action and drawio-export-action are well-maintained, widely adopted, and solve the stated problems correctly. Recommending them is the correct outcome independent of authorship. Suppressing a correct recommendation to avoid the appearance of bias produces a worse outcome than the bias it avoids. The motivation is disclosed. This entry is the disclosure. The recommendation stands.",
		RevisitConditions:   "If a superior alternative is published and maintained to an equivalent standard.",
		MapsTo:              "The recommendation stands. The irony is noted.",
	},
	"DME-3000": {
		ID:                  "DME-3000",
		Title:               "dismissive expletive output mode",
		Status:              "Rejected",
		RejectionConfidence: "100%",
		OptionEvaluated:     "Escalate sigh output beyond MODERATE to include dismissive expletives.",
		RejectionRationale:  "The Exasperation Subsystem's output is capped at MODERATE. Escalation beyond that point was evaluated and rejected. Profanity does not increase signal. It decreases credibility. WRAAS does not decrease credibility.",
		RevisitConditions:   "None. The cap is structural.",
	},
}

// GetDME returns a DME entry by ID.
func GetDME(id string) (DMEEntry, bool) {
	entry, ok := DMEEntries[id]
	return entry, ok
}
