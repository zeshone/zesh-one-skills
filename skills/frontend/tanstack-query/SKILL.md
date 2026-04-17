---
name: tanstack-query
description: >
  ZeshOne TanStack Query v5 conventions.
  Trigger: When managing server/async state — data fetching, caching, mutations, prefetching, infinite scroll.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "1.0"
  inspired-by: frontend-nextjs/tanstack-query.md
---

## When to Use

Load this skill when working with server/async state: data fetching, caching, mutations, optimistic updates, prefetching, or infinite queries.

## Critical Patterns

- Use TanStack Query for server data and keep React Context out of API/cache state.
- Create a shared `QueryClient` with default query and mutation behavior, then expose it through a provider.
- Define feature-level query key factories and reuse them everywhere for fetching, invalidation, and cache updates.
- Wrap fetch definitions in `queryOptions` or `infiniteQueryOptions` so hooks, prefetching, and hydration share the same source of truth.
- Prefetch on the server and hydrate with `HydrationBoundary` when the first render should arrive warm.
- Use optimistic updates carefully: snapshot previous cache, roll back on error, and invalidate the affected detail/list queries on settle.

## State Management Decision Matrix

| Data Type | Solution |
|-----------|----------|
| Server/API data | **TanStack Query** |
| Global UI state | **Zustand** |
| Form state | **React Hook Form** |
| Local component state | **useState** |
| Theme/Locale (rarely-changing) | **React Context** |

**NEVER use React Context for API/server data. ALWAYS use TanStack Query.**

## Setup

```typescript
// lib/query-client.ts
import { QueryClient, isServer } from "@tanstack/react-query";

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000,         // 1 minute
        gcTime: 10 * 60 * 1000,       // 10 minutes
        retry: 1,
        retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 30000),
        refetchOnWindowFocus: false,
      },
      mutations: { retry: 1 },
    },
  });
}

let browserQueryClient: QueryClient | undefined;

export function getQueryClient() {
  if (isServer) return makeQueryClient();
  if (!browserQueryClient) browserQueryClient = makeQueryClient();
  return browserQueryClient;
}
```

```tsx
// providers/query-provider.tsx
"use client";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { getQueryClient } from "@/lib/query-client";

export function QueryProvider({ children }: { children: React.ReactNode }) {
  const queryClient = getQueryClient();
  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
```

## Query Key Factories (ALWAYS Use This Pattern)

```typescript
// features/posts/queries/keys.ts
export const postKeys = {
  all: ["posts"] as const,
  lists: () => [...postKeys.all, "list"] as const,
  list: (filters: PostFilters) => [...postKeys.lists(), { filters }] as const,
  details: () => [...postKeys.all, "detail"] as const,
  detail: (id: string) => [...postKeys.details(), id] as const,
};
```

Key factories ensure consistent cache management and type-safe invalidation across the feature.

## Query Options Pattern

```typescript
// features/posts/queries/options.ts
import { queryOptions, infiniteQueryOptions } from "@tanstack/react-query";
import { postKeys } from "./keys";
import { api } from "@/lib/api/client";

export const postListOptions = (filters: PostFilters) =>
  queryOptions({
    queryKey: postKeys.list(filters),
    queryFn: () => api.posts.list(filters),
  });

export const postDetailOptions = (id: string) =>
  queryOptions({
    queryKey: postKeys.detail(id),
    queryFn: () => api.posts.get(id),
    enabled: !!id,
  });

export const postInfiniteOptions = (filters: PostFilters) =>
  infiniteQueryOptions({
    queryKey: postKeys.list(filters),
    queryFn: ({ pageParam }) => api.posts.list({ ...filters, cursor: pageParam }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.nextCursor ?? undefined,
  });
```

## Server Prefetching with HydrationBoundary

Prefetch in Server Components → client gets data instantly, no loading state on first render.

```tsx
// app/(dashboard)/posts/page.tsx — Server Component
import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { getQueryClient } from "@/lib/query-client";
import { postListOptions } from "@/features/posts/queries/options";

export default async function PostsPage() {
  const queryClient = getQueryClient();
  await queryClient.prefetchQuery(postListOptions({ status: "published" }));
  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <PostList />
    </HydrationBoundary>
  );
}

// features/posts/components/post-list.tsx — Client Component
"use client";
import { useQuery } from "@tanstack/react-query";

export function PostList() {
  // Data immediately available from server prefetch
  const { data: posts } = useQuery(postListOptions({ status: "published" }));
  return (
    <div className="grid gap-4">
      {posts?.map((post) => <PostCard key={post.id} post={post} />)}
    </div>
  );
}
```

## Custom Hooks

```typescript
// features/posts/hooks/use-posts.ts
"use client";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";

export const usePosts = (filters: PostFilters) => useQuery(postListOptions(filters));
export const usePost = (id: string) => useQuery(postDetailOptions(id));
// useSuspenseQuery — use when wrapped in <Suspense> boundary
export const usePostSuspense = (id: string) => useSuspenseQuery(postDetailOptions(id));
```

## Mutations with Optimistic Updates

```typescript
// features/posts/hooks/use-post-mutations.ts
"use client";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useUpdatePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdatePostInput) => api.posts.update(data.id, data),

    onMutate: async (newData) => {
      await queryClient.cancelQueries({ queryKey: postKeys.detail(newData.id) });
      const previous = queryClient.getQueryData(postKeys.detail(newData.id));
      queryClient.setQueryData(postKeys.detail(newData.id), (old: Post) => ({ ...old, ...newData }));
      return { previous };
    },

    onError: (_err, newData, context) => {
      if (context?.previous)
        queryClient.setQueryData(postKeys.detail(newData.id), context.previous);
    },

    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({ queryKey: postKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
    },
  });
}

export function useDeletePost() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.posts.delete(id),
    onSuccess: (_data, id) => {
      queryClient.removeQueries({ queryKey: postKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
    },
  });
}
```

## Infinite Queries (Cursor Pagination)

```tsx
"use client";
import { useInfiniteQuery } from "@tanstack/react-query";
import { useInView } from "react-intersection-observer";

export function PostFeed({ filters }: { filters: PostFilters }) {
  const { ref, inView } = useInView();
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useInfiniteQuery(postInfiniteOptions(filters));

  useEffect(() => {
    if (inView && hasNextPage) fetchNextPage();
  }, [inView, hasNextPage, fetchNextPage]);

  return (
    <div className="flex flex-col gap-4">
      {data?.pages.flatMap((page) =>
        page.items.map((post) => <PostCard key={post.id} post={post} />)
      )}
      <div ref={ref}>{isFetchingNextPage && <Spinner />}</div>
    </div>
  );
}
```

## Cache Invalidation

```typescript
queryClient.invalidateQueries({ queryKey: postKeys.all });           // all posts
queryClient.invalidateQueries({ queryKey: postKeys.lists() });       // only lists
queryClient.invalidateQueries({ queryKey: postKeys.detail("123"), exact: true }); // specific
queryClient.removeQueries({ queryKey: postKeys.detail("123") });     // remove from cache
queryClient.setQueryData(postKeys.detail("123"), updatedPost);       // update directly
```

## Dependent Queries

```typescript
export function usePostWithAuthor(postId: string) {
  const post = useQuery(postDetailOptions(postId));
  const author = useQuery({
    queryKey: userKeys.detail(post.data?.authorId ?? ""),
    queryFn: () => api.users.get(post.data!.authorId),
    enabled: !!post.data?.authorId, // Only fetches when post data is available
  });
  return { post: post.data, author: author.data, isLoading: post.isPending };
}
```

## Resources

- TanStack Query React Overview — https://tanstack.com/query/latest/docs/framework/react/overview
- TanStack Query Query Keys — https://tanstack.com/query/latest/docs/framework/react/guides/query-keys
- TanStack Query Mutations — https://tanstack.com/query/latest/docs/framework/react/guides/mutations
- TanStack Query Infinite Queries — https://tanstack.com/query/latest/docs/framework/react/guides/infinite-queries
- TanStack Query Advanced SSR — https://tanstack.com/query/latest/docs/framework/react/guides/advanced-ssr
