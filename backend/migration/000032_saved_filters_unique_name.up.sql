ALTER TABLE saved_filters ADD CONSTRAINT saved_filters_user_id_name_key UNIQUE (user_id, name);
