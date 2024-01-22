package com.gakaki.demo.service;

import java.util.Collection;
import java.util.List;
import java.util.concurrent.*;
import java.util.stream.Collectors;

public class FixedVirtualThreadExecutorService implements ExecutorService {
    private final ExecutorService VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE = Executors.newVirtualThreadPerTaskExecutor();

    private Semaphore semaphore;

    private int poolSize;

    public FixedVirtualThreadExecutorService(int poolSize) {
        this.poolSize = poolSize;
        this.semaphore = new Semaphore(this.poolSize);
    }

    @Override
    public void shutdown() {
        VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE.shutdown();
    }

    @Override
    public List<Runnable> shutdownNow() {
        return VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE.shutdownNow();
    }

    @Override
    public boolean isShutdown() {
        return VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE.isShutdown();
    }

    @Override
    public boolean isTerminated() {
        return VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE.isTerminated();
    }

    @Override
    public boolean awaitTermination(long timeout, java.util.concurrent.TimeUnit unit) throws InterruptedException {
        return VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE.awaitTermination(timeout, unit);
    }

    @Override
    public <T> Future<T> submit(Callable<T> task) {
        return CompletableFuture.supplyAsync(() -> {
            try {
                semaphore.acquire();
                return task.call();
            } catch (Exception e) {
                throw new IllegalStateException(e);
            } finally {
                semaphore.release();
            }
        }, VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE);

    }

    @Override
    public <T> Future<T> submit(Runnable task, T result) {
        return CompletableFuture.supplyAsync(() -> {
            try {
                semaphore.acquire();
                task.run();
                return result;
            } catch (Exception e) {
                throw new IllegalStateException(e);
            } finally {
                semaphore.release();
            }
        }, VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE);
    }

    @Override
    public Future<?> submit(Runnable task) {
        return CompletableFuture.supplyAsync(() -> {
            try {
                semaphore.acquire();
                task.run();
                return null;
            } catch (Exception e) {
                throw new IllegalStateException(e);
            } finally {
                semaphore.release();
            }
        }, VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE);
    }

    @Override
    public <T> List<Future<T>> invokeAll(Collection<? extends Callable<T>> tasks) throws InterruptedException {
        return tasks.stream().map(t -> CompletableFuture.supplyAsync(() -> {
            try {
                semaphore.acquire();
                return t.call();
            } catch (Exception e) {
                throw new IllegalStateException(e);
            } finally {
                semaphore.release();
            }
        }, VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE)).collect(Collectors.toList());
    }

    @Override
    public <T> List<Future<T>> invokeAll(Collection<? extends Callable<T>> tasks, long timeout, java.util.concurrent.TimeUnit unit) throws InterruptedException {
        return tasks.stream().map(t -> CompletableFuture.supplyAsync(() -> {
            try {
                semaphore.acquire();
                return t.call();
            } catch (Exception e) {
                throw new IllegalStateException(e);
            } finally {
                semaphore.release();
            }
        }, VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE).orTimeout(timeout, unit)).collect(Collectors.toList());
    }

    @Override
    public <T> T invokeAny(Collection<? extends Callable<T>> tasks) throws InterruptedException, ExecutionException {
        return tasks.stream().map(t -> CompletableFuture.supplyAsync(() -> {
            try {
                semaphore.acquire();
                return t.call();
            } catch (Exception e) {
                throw new IllegalStateException(e);
            } finally {
                semaphore.release();
            }
        }, VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE)).map(f -> {
            try {
                return f.get();
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            } catch (ExecutionException e) {
                throw new RuntimeException(e);
            }
        }).findAny().get();
    }

    @Override
    public <T> T invokeAny(Collection<? extends Callable<T>> tasks, long timeout, java.util.concurrent.TimeUnit unit) throws InterruptedException, ExecutionException, TimeoutException {
        return invokeAll(tasks, timeout, unit).stream().map(f -> {
            try {
                return f.get();
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            } catch (ExecutionException e) {
                throw new RuntimeException(e);
            }
        }).findAny().get();
    }

    @Override
    public void close() {
        VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE.close();
    }

    @Override
    public void execute(Runnable command) {
        VIRTUAL_THREAD_POOL_EXECUTOR_SERVICE.execute(() -> {
            try {
                semaphore.acquire();
                command.run();
            } catch (Exception e) {
                throw new IllegalStateException(e);
            } finally {
                semaphore.release();
            }
        });

    }
}