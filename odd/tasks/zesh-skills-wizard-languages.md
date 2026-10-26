# Scope, architecture and language menus
Branch: feat/zesh-skills-wizard-languages; baseline 32ee027200b56cff12ac7fca3786ef8673725df6. Delivery: local work-unit commits only; no source push, PR, release or global installation requested.
Forecast: ~900 authored diff lines. Checks: focused/full Go tests, vet, Node suite, six CGO-disabled builds; independent verification follows native assessment.

## Specs
- S1: "ok quedo bien, agrega la misma seleccion para arquitectura y para scope"
- S2: "adicional quiero que expandas el catalogo de arquitecturas con mas tipos y tipos hibridos comunes, con su breve explicacion"
- S3: "quiero una opcion al inicio de la instalacion o un parametro para elegir entre español o ingles (default ingles) talvez con un parametro basta algo asi\nzesh-skills --es"
- S4: "la CLI debe ser multiplataforma, un binario sin dependencias"
- S5: "las skills solo se instalan localmente en el proyecto a trabajar"

## Tasks
- T1 | S1-S5 | inline exploration/preparation | done | Previous TUI manually accepted; baseline local commit 32ee027; bounded error/protocol mapping read back.
- T2 | S1-S2 | inline (background writers disabled) | done | RED/GREEN scope/profile/public CLI checks, full Go and vet PASS; local commit 9188d03; native medium/reviewDue=false, self verification stands.
- T3 | S3-S5 | inline (same parser/picker unit) | in_progress | Language unit fd897b0 initially GREEN; reopened for independent B1: preserve comma-space in default-English multi-topic --list without changing picker separators.
- T4 | S1-S5 | native-assessment plan + parent checks | pending | Independent report consolidated: one B1 compatibility blocker, no others; after T3 correction, limited independent recheck and refreshed builds required.

## Log
L1 User request (verbatim):
> ok quedo bien, agrega la misma seleccion para arquitectura y para scope
> adicional quiero que expandas el catalogo de arquitecturas con mas tipos y tipos hibridos comunes, con su breve explicacion
> quiero una opcion al inicio de la instalacion o un parametro para elegir entre español o ingles (default ingles) talvez con un parametro basta algo asi
> zesh-skills --es
L2 Previous TUI manual validation closes the prior terminal-exit blocker. Previous work captured locally as 32ee027; no CLI source was published. Catalog main remains 9b099b0. Existing terminal raw/restore adapters and pinned dependencies will not be redesigned.
L3 Technical choices under existing autonomy: reuse a synchronous generic picker; scope accepts several domains, architecture chooses one profile and retains Custom text. Central local architecture catalog includes simple/layered/vertical-slice/Clean/Hexagonal/Onion/modular/distributed/UI styles and common hybrids; descriptions distinguish domain organization, deployment styles and patterns. Selection stores a stable profile ID, not a translated label, and creates no infrastructure or mandatory skill bundle.
L4 --es is the sufficient requested language mechanism; no extra startup question or persisted locale. English remains byte-compatible by default. Localize owned prompts, menus, help, preview, confirmations and diagnostics at presentation boundaries. Do not translate IDs, flags, JSON keys, hashes, custom user values, downloaded skill documents, original catalog descriptions or generated agent instruction formats. OS/technical details may retain their original language, explicitly identified in Spanish.
L5 Explorer evidence: English errors also originate in catalog.go/install.go, including a modified-index confirmation suffix. Existing verification_test.go asserts exact default-English preview/error streams and a three-read conflict fixture. Preserve these by leaving domain/storage code unchanged and localizing UI templates before inserting opaque values.
L10 Independent verifier mvoxqov1-8-wxni: Go/vet/Node PASS (known symlink SKIP), but HOLD B1 confirmed by baseline source diff: --list formerly joined topics with comma-space, now language helper joins comma only. Existing single-topic fixtures miss it. No executable reproduction claimed yet. T3 reopened for a public RED/GREEN regression and one bounded correction; T4 awaits targeted recheck. Other S1/S2/S4/S5 and Spanish behavior pass.
L9 Final artifact evidence: job10 six CGO-disabled builds PASS; job11 PASS for six hashes, English/default and --es help/errors/list/install, --es=false, empty PATH, real main catalog9b099b0, stable architecture ID, preserved CRLF/index and idempotent English repeat. Retained project dist/wizard-language-smoke-S2VSOY/project. Closing assessment over base32ee027 escalated shared parser/picker contracts (item3): nativeRisk medium, risk high for verification, reviewDue=false, independentVerifier=true. One independent verification now required; no native review or old winpty retry.
L8 T3 evidence: RED all Spanish tests rejected unknown --es, then full Go/vet/Node GREEN. Public checks cover help/errors without network, Spanish custom/hybrid menus, sí confirmation, preserved English errors/previews and byte-identical stored files, optional --es=false, cancel/conflict preservation. Local commit fd897b0; native medium372lines/reviewDue=false/independentVerifier=false. Final six-target rebuild job10 started.
L7 T2 evidence: RED undefined architectureProfiles/newOptionPicker/scopeOptions, then full Go/focused menu checks/vet GREEN. Twenty bilingual profile definitions include UI/distributed styles and five common hybrids; stable IDs/custom text checked. Local commit 9188d03; native committed-range assessment medium,340 lines, reviewDue=false, independentVerifier=false. No native review started for this unit.
L6 Test-first: each new menu/language rule gets a RED public-interface/model test; changed shared picker must retain old skill key behavior and automation flags. Scope/architecture cancellation must not fetch or write; --help/--es --help need no network. Native assessment decides outside verification/review, never a TODO transition.
