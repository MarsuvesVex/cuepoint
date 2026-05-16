export type TwitchHelixResponse<T> = {
  data: T[];
  pagination?: {
    cursor?: string;
  };
};

export type TwitchUserProfile = {
  id: string;
  login: string;
  display_name: string;
};

export type TwitchChannel = {
  broadcaster_id: string;
  broadcaster_login: string;
  broadcaster_name: string;
  title: string;
};

export type TwitchStream = {
  id: string;
  title: string;
  started_at: string;
};

export type TwitchStreamMarker = {
  id: string;
  created_at: string;
  description: string;
  position_seconds: number;
};

export type TwitchVideo = {
  id: string;
  stream_id: string;
  user_id: string;
  user_login: string;
  user_name: string;
  title: string;
  created_at: string;
  url: string;
  duration: string;
};

export type TwitchSchedule = {
  broadcaster_id: string;
  broadcaster_name: string;
  broadcaster_login: string;
  segments: TwitchScheduleSegment[];
};

export type TwitchScheduleSegment = {
  id: string;
  start_time: string;
  end_time: string;
  title: string;
  is_recurring: boolean;
  category: {
    id: string;
    name: string;
  } | null;
};
