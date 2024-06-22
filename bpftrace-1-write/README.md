Однажды вы заходите на сервер и понимаете, что кто-то очень много пишет на диск. С помощью стандартных утилит (кстати, каких?) вы понимаете, что дело в вашем микросервисе. Вы открываете код микросервиса, и понимаете, что запись на диск идёт из многих разных мест и конечно-же ничего не логируется. Патчить код, чтобы добавить логирование, выкатывать новую версию - слишком долго, к тому же потом придётся писать скрипты для анализа логов. Вы вспоминаете, что есть утилита `bpftrace`, которая как раз подходит для таких случаев!

Напишите трейсфайл `trace.bt` для bpftrace, который бы распечатал бектрейсы, приводящие к записи на диск и статистику переданных байт.

Для примера, попробуйте потрейсить микросервис `main.go` (можно считать, что `main` - имя бираника микросервиса). Должно получиться что-то вроде:


```
@stacks[
    syscall.Syscall.abi0+27
    internal/poll.(*FD).Write+878
    os.(*File).Write+101
    main.write+200
    main.background+40
    main.main.func1+46
    runtime.goexit.abi0+1
]: 8192
@stacks[
    syscall.Syscall.abi0+27
    internal/poll.(*FD).Write+878
    os.(*File).Write+101
    main.write+200
    main.background+40
    main.main.func1+46
    runtime.goexit.abi0+1
]: 12288
@stacks[
    syscall.Syscall.abi0+27
    internal/poll.(*FD).Write+878
    os.(*File).Write+101
    main.write+200
    main.foreground+40
    main.main+263
    runtime.main+530
    runtime.goexit.abi0+1
]: 16384
@stacks[
    syscall.Syscall.abi0+27
    internal/poll.(*FD).Write+878
    os.(*File).Write+101
    main.write+200
    main.foreground+40
    main.main+263
    runtime.main+530
    runtime.goexit.abi0+1
]: 20480
```

Горутины могут засыпать, что приводит к появлению трейсов, содержащую функцию `time.Sleep`. Такие трейсы можно игнорировать. Например, такой бектрейс в выводе `bpftrace` будет проигнорирован чекером задания:

```
@stacks[
    runtime.write1.abi0+21
    runtime.netpollBreak+88
    runtime.wakeNetPoller+52
    runtime.modtimer+968
    runtime.resetForSleep+55
    runtime.park_m+167
    runtime.mcall+67
    time.Sleep+302
    main.foreground+50
    main.main+263
    runtime.main+530
    runtime.goexit.abi0+1
]: 2
```
