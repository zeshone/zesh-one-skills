using System.Globalization;

namespace CliAcceptance;

internal static class Program
{
    private const string UsageError = "invalid arguments. See --help.";
    private const string Help =
        "Usage: line-tally [--json] [--] <file|->\n" +
        "Count text lines without changing input.\n\n" +
        "Options:\n" +
        "  --help  Show help and exit.\n" +
        "  --json  Emit {\"lines\":N} instead of a decimal count.\n" +
        "  --      Treat remaining arguments as operands.\n\n" +
        "Input: use - to read stdin; exactly one input is required.\n" +
        "Exit codes: 0 success/help, 2 invalid arguments, 1 input read failure.\n" +
        "Example: line-tally --json input.txt\n";

    private static Task<int> Main(string[] args) => RunAsync(args, CancellationToken.None);

    private static async Task<int> RunAsync(string[] args, CancellationToken cancellationToken)
    {
        if (args is ["--help"])
        {
            await Console.Out.WriteAsync(Help.AsMemory(), cancellationToken);
            return 0;
        }

        bool json = false;
        bool options = true;
        string? input = null;
        foreach (string arg in args)
        {
            if (options && arg == "--") { options = false; continue; }
            if (options && arg == "--json" && !json) { json = true; continue; }
            if (string.IsNullOrWhiteSpace(arg) || input is not null ||
                (options && arg.StartsWith('-') && arg != "-"))
            {
                return await FailAsync(UsageError, 2, cancellationToken);
            }
            input = arg;
        }
        if (input is null) return await FailAsync(UsageError, 2, cancellationToken);

        long count = 0;
        try
        {
            using TextReader reader = input == "-" ? Console.In : File.OpenText(input);
            while (await reader.ReadLineAsync(cancellationToken) is not null) count++;
        }
        catch (Exception error) when (error is IOException or UnauthorizedAccessException)
        {
            return await FailAsync("unable to read input.", 1, cancellationToken);
        }

        // Only invariant digits enter this fixed JSON shape, never user text.
        string digits = count.ToString(CultureInfo.InvariantCulture);
        string output = json ? $"{{\"lines\":{digits}}}\n" : $"{digits}\n";
        await Console.Out.WriteAsync(output.AsMemory(), cancellationToken);
        return 0;
    }

    private static async Task<int> FailAsync(string message, int code, CancellationToken cancellationToken)
    {
        await Console.Error.WriteAsync($"line-tally: {message}\n".AsMemory(), cancellationToken);
        return code;
    }
}
