using SyncClipboard.Core.Clipboard;
using SyncClipboard.Core.Interfaces;
using SyncClipboard.Core.Models;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace SyncClipboard.Desktop.ClipboardAva;

internal class ClipboardListener(IClipboardFactory clipboardFactory, ILogger logger) : ClipboardChangingListenerBase
{
    protected override IClipboardFactory ClipboardFactory { get; } = clipboardFactory;
    private readonly ILogger _logger = logger;
    // 获取电脑名字
    private readonly string _source = Environment.MachineName;

    private Timer? _timer;
    private MetaChanged? _action;
    private ClipboardMetaInfomation? _meta;

    private readonly SemaphoreSlim _tickSemaphore = new(1, 1);
    private CancellationTokenSource? _cts;

    protected override void RegistSystemEvent(MetaChanged action)
    {
        _action = action;
        _timer = new Timer(InvokeTick, null, TimeSpan.Zero, TimeSpan.FromSeconds(1));
    }

    protected override void UnRegistSystemEvent(MetaChanged action)
    {
        _timer?.Dispose();
        _timer = null;

        _action = null;

        _cts?.Cancel();
        _cts?.Dispose();
        _cts = null;
    }

    internal void TriggerClipboardChangedEvent()
    {
        _cts?.Cancel();
        _cts?.Dispose();
        _cts = null;
        InvokeTick(null);
    }

    private async void InvokeTick(object? _)
    {
        if (_tickSemaphore.Wait(0) is false)
        {
            return;
        }

        try
        {
            _cts?.Dispose();
            _cts = new CancellationTokenSource(TimeSpan.FromSeconds(10));

            var meta = await ClipboardFactory.GetMetaInfomation(_cts.Token);

            // 使用 with 表达式创建新对象，同时设置 Source
            var metaWithSource = meta with { Source = _source };

            if (metaWithSource == _meta)
            {
                return;
            }

            if (_meta is not null)
            {
                _meta = metaWithSource;
                _ = Task.Run(() => _action?.Invoke(metaWithSource));
                _ = _logger.WriteAsync($"Clipboard changed to {metaWithSource}");
            }
            else
            {
                _meta = metaWithSource;
            }
        }
        catch { }
        finally
        {
            _tickSemaphore.Release();
        }
    }
}
