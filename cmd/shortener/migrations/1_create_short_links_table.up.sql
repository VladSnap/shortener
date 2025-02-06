CREATE TABLE IF NOT EXISTS public.short_links (
  uuid UUID NOT NULL,
  short_url varchar NOT NULL,
  orig_url varchar NOT NULL,
  PRIMARY KEY (uuid)
);
CREATE UNIQUE INDEX IF NOT EXISTS short_links_short_url_unique_idx on public.short_links (short_url);
CREATE UNIQUE INDEX IF NOT EXISTS short_links_orig_url_unique_idx on public.short_links (orig_url);