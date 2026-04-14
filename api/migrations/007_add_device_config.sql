alter table devices add column if not exists config jsonb not null default '{"min_contour_area": 500, "threshold": 25, "cooldown_seconds": 2.0}';
