# Why pomo works the way it does

pomo is a focus timer, not a treatment. The design is an attempt to use the
parts of the ADHD literature that have the best support, and to avoid the parts
that are popular but untested. This document is deliberately blunt about the
difference.

If you take one thing from it: **the timer is a container, not the
intervention.** The evidence is strongest for structured time management,
external time cues, one clear next action, immediate feedback, and offloading
intrusive thoughts. It is weakest for the specific 25/5 protocol, “dopamine”
rewards, body doubling, and “hyperfocus mode.”

## Evidence grades

| Grade | Meaning |
|---|---|
| **A** | Meta-analysis or RCT in **adults with ADHD** for this kind of intervention |
| **B** | RCT in ADHD, but pediatric, multimodal, or otherwise indirect for an adult using a timer |
| **C** | General-population or mechanistic evidence; not directly tested in adult ADHD |
| **D** | No direct evidence; clinical rationale, user autonomy, or harm reduction only |

Grades describe the **intervention component**, not this app. There is no RCT of
pomo itself.

## What the evidence supports, and how pomo uses it

### 1. External time is compensatory, not motivational

ADHD is associated with timing deficits across motor timing, duration
discrimination, and temporal foresight, with dysfunction in fronto-striatal and
fronto-parietal timing networks (Noreika et al., 2013). This is why “just feel
the time” fails: the internal clock is not a reliable instrument.

The best direct intervention evidence is a randomized controlled trial in
children and adolescents with ADHD (n=38 analyzed). A multimodal intervention
combining time-skill training with time-assistive devices — including visual
timers — improved parent-rated daily time management (d=1.0) and
time-processing ability (d=.38). The children’s own time-management ratings did
not improve significantly, and the timer was not tested alone (Wennberg et al.,
2018).

Adult evidence is for **time-management training**, not a visual timer in
isolation:

- Meta-cognitive therapy focused on time management, organization, and planning
  beat supportive psychotherapy in 88 adults (Solanto et al., 2010; response
  OR 5.41).
- A pilot RCT of group CBT focused on time management (n=48) reduced ADHD
  symptoms on all measures (Nakashima et al., 2022).
- CBT for adult ADHD has small-to-medium effects versus control (g≈.65 for
  symptoms, g≈.51 for functioning; Knouse et al., 2017).

**In pomo:** a large countdown, a shrinking progress bar, and “ends at 14:32.”
This is an external time cue, not a claim that a timer treats ADHD.

### 2. Start with one small, concrete action

Behavioral activation targets the initiation problem directly. In a 2025 RCT
of adults with ADHD predominantly inattentive presentation (n=108), a CBT
protocol that included behavioral activation improved activation (d=.49;
activation subscale d=.65) more than an active CBT comparator. Between-group
procrastination did not improve significantly, and the trial was underpowered
(Strålin et al., 2025).

A shorter CBT planning intervention for college students with ADHD (n=58)
improved inattention symptoms but was modest and did not separate on most
outcomes (Van der Oord et al., 2020). The takeaway is not “all planning works.”
It is that **graded activation plus concrete planning is a reasonable target**,
and that effects are modest.

**In pomo:** a 5-minute activation sprint (`1`), not the folk “2-minute rule.”
The exact 5 minutes is a design choice; the evidence supports graded activation,
not a specific number. After the sprint, the app offers to continue rather than
forcing a break.

### 3. One next action beats a backlog

Organization, planning, and time management are core components of the adult
ADHD CBT protocols with RCT support (Solanto et al., 2010; Knouse et al., 2017).
The CADDI protocol explicitly trains dividing big tasks into smaller pieces and
defining actions that fit goals (Strålin et al., 2025). Proximal subgoals also
have a long history in goal-setting and procrastination research, though not
specifically in adult ADHD.

**In pomo:** one line for the next physical action, always visible. There is no
project tree, backlog, or priority matrix. Press `i` to change it.

### 4. Immediate feedback, not distant points

Adults with ADHD show steeper monetary delay discounting than controls
(d=.40 in adult samples; d=.43 overall; Jackson & MacKillop, 2016). Children
with ADHD respond well to immediate, high-intensity reinforcement, and prefer
immediate over delayed rewards; the strongest evidence is in children (Luman
et al., 2005 review).

This supports **immediate, salient progress**, not a slot machine. It does not
prove that confetti, points, or variable rewards work in adults.

**In pomo:** at the end of a block, the screen immediately shows
`+25m banked` and today’s total. No coins, no levels, no delayed streak to
protect.

### 5. Capture distractions so they stop consuming working memory

Adults with ADHD show planning impairments on a complex prospective-memory
paradigm, even when execution of already-formed intentions is relatively intact
(Fuermaier et al., 2013). A meta-analysis of prospective-memory interventions
found external memory aids to be the most effective approach, while noting that
larger trials are needed (Jones et al., 2021). Cognitive offloading — writing
something down, setting a reminder — reduces internal cognitive demand (Risko &
Gilbert, 2016).

**In pomo:** `d` opens a one-line capture box and returns to the timer. The
thought is stored in the parking lot and offered for review at the break. The
app does not claim this will eliminate distraction; it removes the need to hold
the thought in working memory.

### 6. If-then plans, with a caveat

Implementation intentions (“if situation X, then I will do Y”) have medium-to-
large effects on goal attainment in the general population (Gollwitzer &
Sheeran, 2006). In ADHD, the evidence is mostly from children: if-then plans
improved response inhibition and delay of gratification (Gawrilow et al.,
2011).

I found **no adult ADHD RCT** of implementation intentions for focus sessions.
Treat this as promising and cheap, not proven for this population.

**In pomo:** an optional second line in the intention editor, e.g.
`if I want to check email, then I'll park it`. It is optional because it adds
setup friction, and the adult evidence is thin.

### 7. Breaks are for fatigue, not a productivity cheat code

A meta-analysis of micro-breaks (22 samples, n=2335) found small benefits for
vigor (d=.36) and fatigue (d=.35), but **no significant overall performance
effect** (d=.16). Performance benefits appeared only for less cognitively
demanding tasks, and longer breaks helped more after demanding work (Albulescu
et al., 2022).

Studies of Pomodoro-style fixed breaks are mixed:

- One study found mood and efficiency benefits versus self-regulated breaks,
  with no difference in task completion (Biwer et al., 2023).
- A later 2-hour study found no overall differences in motivation, fatigue,
  productivity, task completion, or flow; Pomodoro and Flowtime showed a faster
  decline in motivation over time (Smits et al., 2025).
- A 2025 scoping review of the Pomodoro technique found 32 studies (total
  N=5,270), but only 3 RCTs (n=87), mostly in education; the authors note
  methodological diversity and limited comparability, and call for long-term
  trials (Ogut, 2025).

**In pomo:** breaks are skippable, short by default, and framed as fatigue
management. The break screen suggests standing, water, or looking outside —
not because “screen-free” is proven, but because it is a low-cost alternative
to scrolling. Auto-start break is on; auto-start focus is off, preserving
agency at the harder transition.

### 8. One boundary cue, not a reminder system

A randomized trial of SMS reminders in adults with ADHD found no improvement in
module completion or number of logins; reminders only produced faster login in
some modules (Nordby et al., 2022). More notification is not the answer.

**In pomo:** one bell and one optional desktop notification at the end of a
block. No repeated “are you working?” pings.

### 9. Passive self-monitoring

Self-monitoring interventions improve on-task behavior in students with ADHD,
with the strongest evidence from single-case designs and student samples. Adult
evidence is thinner. Mandatory ratings also add friction and can become a
guilt ritual.

**In pomo:** a passive session log and a neutral 7-day view. No end-of-block
interrogation, no completion percentage, no red failures.

### 10. Extending a block is autonomy, not “hyperfocus mode”

“Hyperfocus” is poorly defined, has no consensus definition, and has almost no
intervention research (Ashinoff & Abu-Akel, 2021). Adults with more ADHD
symptoms report more hyperfocus on a questionnaire (Hupfeld et al., 2019), but
that does not make it a therapeutic target.

**In pomo:** `e` adds 10 minutes because abruptly interrupting deep engagement
is aversive and because adults should control their own session. It is not
marketed as a treatment, and the app gently offers a break after repeated
extensions.

## What pomo deliberately does **not** claim

- **Body doubling.** No RCT. A 2025 VR preprint (n=12) and an industry survey
  are not enough. Excluded from v1. The preprint itself notes that prior work
  found no conclusive effectiveness.
- **“Dopamine rewards.”** Reinforcement matters, but “dopamine hacking” is
  pop neuroscience. The app uses immediate progress feedback, not dopamine
  claims.
- **25/5 is optimal.** It is a familiar default. Every duration is editable.
- **Streaks build habits.** No adult-ADHD evidence that streak loss helps;
  it can add task aversiveness. pomo shows neutral totals.
- **Phone presence destroys focus.** The “brain drain” finding failed direct
  replication. Break suggestions are a preference, not a proven rule.
- **Notifications fix follow-through.** See Nordby et al., 2022.
- **This app replaces medication or therapy.** It does not. Medication and
  structured psychological treatment are the evidence-based treatments; this
  is a small compensatory tool.

## Other design choices that are **not** evidence claims

- **Tokyo Night palette, big block digits, rounded panels.** Visual salience
  and legibility. No claim that color treats ADHD.
- **8-second auto-break countdown.** A short decision window that still lets
  you extend or skip. The exact number is arbitrary.
- **No settings screen in v1.** Fewer decisions at launch. Config is a JSON
  file for people who want it.
- **One tick per second, none while paused.** Battery and CPU discipline, not
  a clinical claim.
- **Single-instance lock, atomic writes.** Prevent data loss, not a clinical
  claim.

## Limitations

- Most ADHD-specific RCTs are in children and adolescents; adult evidence is
  mostly CBT and time-management training, not timer apps.
- Effect sizes are group averages. They do not predict your response.
- Several apparent mechanisms are indirect: external memory aids are not
  ADHD-specific, and implementation-intention evidence in adults with ADHD is
  missing.
- The app has no RCT. It should be judged by whether it reduces friction for
  you, not by the evidence table.
- The stats measure **completed timer blocks**, not verified work.

## References

Albulescu, P., Macsinga, I., Rusu, A., Sulea, C., Bodnaru, A., & Tulbure, B. T.
(2022). “Give me a break!” A systematic review and meta-analysis on the efficacy
of micro-breaks for increasing well-being and performance. *PLoS ONE*, 17(8),
e0272460. https://doi.org/10.1371/journal.pone.0272460

Ara, Z., Rahim, I. B., Zhou, P., Yu, L., Esmaeili, B., Yu, L.-F., & Hong, S. R.
(2025). You are not alone: Designing body doubling for ADHD in virtual reality.
*arXiv preprint* arXiv:2509.12153. https://arxiv.org/abs/2509.12153

Ashinoff, B. K., & Abu-Akel, A. (2021). Hyperfocus: the forgotten frontier of
attention. *Psychological Research*, 85(1), 1–19.
https://doi.org/10.1007/s00426-019-01245-8

Biwer, F., Wiradhany, W., Oude Egbrink, M. G. A., & de Bruin, A. B. H. (2023).
Understanding effort regulation: Comparing “Pomodoro” breaks and self-regulated
breaks. *British Journal of Educational Psychology*, 93(S2), 353–367.
https://doi.org/10.1111/bjep.12593

Fuermaier, A. B. M., Tucha, L., Koerts, J., Aschenbrenner, S., Westermann, C.,
Weisbrod, M., Lange, K. W., & Tucha, O. (2013). Complex prospective memory in
adults with attention deficit hyperactivity disorder. *PLoS ONE*, 8(3), e58338.
https://doi.org/10.1371/journal.pone.0058338

Gawrilow, C., Gollwitzer, P. M., & Oettingen, G. (2011). If-then plans benefit
delay of gratification performance in children with and without ADHD.
*Cognitive Therapy and Research*, 35, 442–455.
https://doi.org/10.1007/s10608-010-9309-z

Gollwitzer, P. M., & Sheeran, P. (2006). Implementation intentions and goal
achievement: A meta-analysis of effects and processes. *Advances in Experimental
Social Psychology*, 38, 69–119.
https://doi.org/10.1016/S0065-2601(06)38002-1

Hupfeld, K. E., Abagis, T. R., & Shah, P. (2019). Living “in the zone”:
hyperfocus in adult ADHD. *ADHD Attention Deficit and Hyperactivity Disorders*,
11(2), 191–208. https://doi.org/10.1007/s12402-018-0272-y

Jackson, J. N. S., & MacKillop, J. (2016). Attention-deficit/hyperactivity
disorder and monetary delay discounting: A meta-analysis of case-control
studies. *Biological Psychiatry: Cognitive Neuroscience and Neuroimaging*, 1(4),
316–325. https://doi.org/10.1016/j.bpsc.2016.01.007

Jones, W. E., Benge, J. F., & Scullin, M. K. (2021). Preserving prospective
memory in daily life: A systematic review and meta-analysis of mnemonic
strategy, cognitive training, external memory aid, and combination
interventions. *Neuropsychology*, 35(1), 123–140.

Knouse, L. E., Teller, J., & Brooks, M. A. (2017). Meta-analysis of
cognitive-behavioral treatments for adult ADHD. *Journal of Consulting and
Clinical Psychology*, 85(7), 737–750.

Luman, M., Oosterlaan, J., & Sergeant, J. A. (2005). The impact of reinforcement
contingencies on AD/HD: A review and theoretical appraisal. *Clinical Psychology
Review*, 25(2), 183–213. https://doi.org/10.1016/j.cpr.2004.11.001

Nakashima, M., Inada, N., Tanigawa, Y., Yamashita, M., Maeda, E., Kouguchi, M.,
et al. (2022). Efficacy of group cognitive behavior therapy targeting time
management for adults with attention deficit/hyperactivity disorder in Japan:
A randomized control pilot trial. *Journal of Attention Disorders*, 26(3),
377–390. https://doi.org/10.1177/1087054720986939

Noreika, V., Falter, C. M., & Rubia, K. (2013). Timing deficits in
attention-deficit/hyperactivity disorder (ADHD): Evidence from neurocognitive
and neuroimaging studies. *Neuropsychologia*, 51(2), 235–266.
https://doi.org/10.1016/j.neuropsychologia.2012.09.036

Nordby, E. S., Gjestad, R., Kenter, R. M. F., Guribye, F., Mukhiya, S. K.,
Lundervold, A. J., & Nordgreen, T. (2022). The effect of SMS reminders on
adherence in a self-guided internet-delivered intervention for adults with ADHD.
*Frontiers in Digital Health*, 4, 821031.
https://doi.org/10.3389/fdgth.2022.821031

Ogut, E. (2025). Assessing the efficacy of the Pomodoro technique in enhancing
anatomy lesson retention during study sessions: A scoping review. *BMC Medical
Education*, 25, 1440. https://doi.org/10.1186/s12909-025-08001-0

Risko, E. F., & Gilbert, S. J. (2016). Cognitive offloading. *Trends in
Cognitive Sciences*, 20(9), 676–688. https://doi.org/10.1016/j.tics.2016.07.002

Smits, E. J. C., Wenzel, N., & de Bruin, A. (2025). Investigating the
effectiveness of self-regulated, Pomodoro, and Flowtime break-taking techniques
among students. *Behavioral Sciences*, 15(7), 861.
https://doi.org/10.3390/bs15070861

Solanto, M. V., Marks, D. J., Wasserstein, J., Mitchell, K., Abikoff, H.,
Alvir, J. M., & Kofman, M. D. (2010). Efficacy of meta-cognitive therapy for
adult ADHD. *American Journal of Psychiatry*, 167(8), 958–968.
https://doi.org/10.1176/appi.ajp.2009.09081123

Strålin, E. E., Thorell, L. B., Lundgren, T., Bölte, S., & Bohman, B. (2025).
Cognitive behavioral therapy for ADHD predominantly inattentive presentation:
Randomized controlled trial of two psychological treatments. *Frontiers in
Psychiatry*, 16, 1564506. https://doi.org/10.3389/fpsyt.2025.1564506

Van der Oord, S., Boyer, B. E., Van Dyck, L., Mackay, K. J., De Meyer, H., &
Baeyens, D. (2020). A randomized controlled study of a cognitive behavioral
planning intervention for college students with ADHD: An effectiveness study in
student counseling services in Flanders. *Journal of Attention Disorders*,
24(6), 849–862. https://doi.org/10.1177/1087054718787033

Wennberg, B., Janeslätt, G., Kjellberg, A., & Gustafsson, P. A. (2018).
Effectiveness of time-related interventions in children with ADHD aged 9–15
years: A randomized controlled study. *European Child & Adolescent Psychiatry*,
27(3), 329–342. https://doi.org/10.1007/s00787-017-1052-5
