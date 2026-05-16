import { notFound, redirect } from "next/navigation";

import { AppShell } from "@/components/app-shell";
import { FormField } from "@/components/form-field";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { addSegmentAction, savePlanAction } from "@/lib/actions";
import { requireSession } from "@/lib/auth";
import { getPlan } from "@/lib/data";

export const dynamic = "force-dynamic";

export default async function PlanDetailPage({
  params,
}: {
  params: Promise<{ planId: string }>;
}) {
  try {
    const session = await requireSession();
    const { planId } = await params;
    const data = await getPlan(planId, session.user.id);

    if (!data.plan) {
      notFound();
    }

    return (
      <AppShell>
        <div className="grid gap-6 xl:grid-cols-[0.95fr_1.05fr]">
          <Card className="space-y-6">
            <div>
              <p className="text-xs uppercase tracking-[0.2em] text-black/45 dark:text-white/45">
                Plan detail
              </p>
              <h2 className="text-2xl font-semibold">{data.plan.name}</h2>
              <p className="mt-2 text-sm text-black/60 dark:text-white/60">
                Define the stream window and the defaults that each segment sequence should inherit.
              </p>
            </div>
            <form action={savePlanAction} className="grid gap-4">
              <input type="hidden" name="plan_id" value={data.plan.id} />
              <FormField
                label="Name"
                hint="A concise show name that appears throughout the control center."
              >
                <Input name="name" defaultValue={data.plan.name} />
              </FormField>
              <FormField
                label="Date"
                hint="The calendar day this plan should appear under."
              >
                <Input
                  name="plan_date"
                  type="date"
                  required
                  defaultValue={String(data.plan.plan_date).slice(0, 10)}
                />
              </FormField>
              <FormField
                label="Timezone"
                hint="Clock-based segments use this timezone when evaluating their start time."
              >
                <Input name="timezone" defaultValue={data.plan.timezone} />
              </FormField>
              <FormField
                label="Planned start"
                hint="Use this when you want the plan itself to carry the day’s expected stream start."
              >
                <Input
                  name="planned_start_at"
                  type="datetime-local"
                  defaultValue={
                    data.plan.planned_start_at
                      ? new Date(data.plan.planned_start_at).toISOString().slice(0, 16)
                      : ""
                  }
                />
              </FormField>
              <FormField
                label="Status"
                hint="Keep draft while editing, scheduled when ready, and live only when the runtime is actively following it."
              >
                <Select name="status" defaultValue={data.plan.status}>
                  <option value="draft">Draft</option>
                  <option value="scheduled">Scheduled</option>
                  <option value="live">Live</option>
                  <option value="completed">Completed</option>
                  <option value="archived">Archived</option>
                </Select>
              </FormField>
              <label className="flex items-center gap-3">
                <Switch
                  name="auto_title_enabled"
                  defaultChecked={data.plan.auto_title_enabled}
                />
                <div>
                  <span className="text-sm font-medium">Automatic title updates</span>
                  <p className="text-xs text-black/55 dark:text-white/55">
                    When enabled, runtime automation may push segment-driven Twitch title changes for this plan.
                  </p>
                </div>
              </label>
              <Button type="submit">Save plan</Button>
            </form>
          </Card>

          <Card className="space-y-6">
            <div>
              <p className="text-xs uppercase tracking-[0.2em] text-black/45 dark:text-white/45">
                Segments
              </p>
              <h3 className="text-2xl font-semibold">Ordered block list</h3>
              <p className="mt-2 text-sm text-black/60 dark:text-white/60">
                Mix clock-based segments with elapsed-time segments. Each block can later feed overlays, title changes, and markers.
              </p>
            </div>
            <div className="space-y-3">
              {data.segments.length === 0 ? (
                <p className="text-sm text-black/55 dark:text-white/55">No segments yet.</p>
              ) : (
                data.segments.map((segment) => (
                  <div
                    key={segment.id}
                    className="rounded-2xl border border-black/10 bg-sand/60 px-4 py-3 dark:border-white/10 dark:bg-white/5"
                  >
                    <div className="flex items-center justify-between gap-4">
                      <div>
                        <p className="font-medium">
                          {segment.position}. {segment.segment_title}
                        </p>
                        <p className="text-sm text-black/55 dark:text-white/55">
                          {segment.segment_type} • {segment.timing_mode} • {segment.duration_estimate_seconds}s
                        </p>
                      </div>
                      <span className="text-xs uppercase tracking-[0.18em] text-black/40 dark:text-white/40">
                        {segment.metadata_status}
                      </span>
                    </div>
                    {segment.youtube_video_title ? (
                      <p className="mt-2 text-sm text-black/60 dark:text-white/60">
                        {segment.youtube_video_title} • {segment.youtube_creator_name}
                      </p>
                    ) : null}
                  </div>
                ))
              )}
            </div>
            <form action={addSegmentAction} className="grid gap-4 rounded-2xl bg-sand/50 p-4 dark:bg-white/5">
              <input type="hidden" name="plan_id" value={data.plan.id} />
              <div className="grid gap-4 md:grid-cols-2">
                <FormField
                  label="Segment title"
                  hint="This is the operator-facing name of the block and the fallback segment title for automation."
                  className="md:col-span-2"
                >
                  <Input name="segment_title" required />
                </FormField>
                <FormField
                  label="Segment type"
                  hint="Standard is a normal block. React enables YouTube metadata and related commands."
                >
                  <Select name="segment_type" defaultValue="standard">
                    <option value="standard">Standard</option>
                    <option value="react">React</option>
                  </Select>
                </FormField>
                <FormField
                  label="Timing mode"
                  hint="Clock uses an absolute time. Elapsed uses time since the stream started."
                >
                  <Select name="timing_mode" defaultValue="elapsed">
                    <option value="elapsed">Elapsed stream time</option>
                    <option value="clock">Clock time</option>
                  </Select>
                </FormField>
                <FormField
                  label="Start time"
                  hint="Use this when timing mode is clock. It represents the wall-clock start for the segment."
                >
                  <Input name="start_at_local" type="datetime-local" />
                </FormField>
                <FormField
                  label="Elapsed offset (seconds)"
                  hint="Use this when timing mode is elapsed. For example, 1800 means the segment should activate 30 minutes after going live."
                >
                  <Input name="elapsed_offset_seconds" type="number" />
                </FormField>
                <FormField
                  label="Duration estimate (seconds)"
                  hint="This is the planned runtime of the segment and helps determine when the next segment should take over."
                >
                  <Input name="duration_estimate_seconds" type="number" defaultValue="0" />
                </FormField>
                <FormField
                  label="Title override"
                  hint="Optional Twitch title text for this segment. Leave blank to use the segment title instead."
                >
                  <Input name="stream_title_override" />
                </FormField>
                <FormField
                  label="YouTube URL"
                  hint="For react segments, Cuepoint queues metadata enrichment and overlay text from this URL."
                  className="md:col-span-2"
                >
                  <Input name="youtube_url" placeholder="https://www.youtube.com/watch?v=..." />
                </FormField>
                <FormField
                  label="Notes"
                  hint="Operator-only planning notes, talking points, or checkpoint reminders."
                  className="md:col-span-2"
                >
                  <Textarea name="notes" />
                </FormField>
                <label className="flex items-center gap-3">
                  <Switch name="react_start_command_enabled" defaultChecked />
                  <div>
                    <span className="text-sm">Enable `!react`</span>
                    <p className="text-xs text-black/55 dark:text-white/55">
                      Allows chat to kick off the active react segment behavior.
                    </p>
                  </div>
                </label>
                <label className="flex items-center gap-3">
                  <Switch name="react_title_command_enabled" defaultChecked />
                  <div>
                    <span className="text-sm">React title updates</span>
                    <p className="text-xs text-black/55 dark:text-white/55">
                      Lets react flows update the Twitch stream title using video metadata.
                    </p>
                  </div>
                </label>
                <label className="flex items-center gap-3 md:col-span-2">
                  <Switch name="react_watching_command_enabled" defaultChecked />
                  <div>
                    <span className="text-sm">Enable `!watching`</span>
                    <p className="text-xs text-black/55 dark:text-white/55">
                      Allows chat to request the active react video title, creator, and canonical link.
                    </p>
                  </div>
                </label>
              </div>
              <Button type="submit">Add segment</Button>
            </form>
          </Card>
        </div>
      </AppShell>
    );
  } catch {
    redirect("/login");
  }
}
