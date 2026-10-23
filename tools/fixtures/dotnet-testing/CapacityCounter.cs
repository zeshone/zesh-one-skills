namespace DotnetTestingTrial;

// Small test specimen: fixed capacity, observable total, no external infrastructure.
public sealed class CapacityCounter
{
    public const int Capacity = 10;
    public int Total { get; private set; } = 8;

    public bool TryAdd(int amount)
    {
        if (amount <= 0 || amount > Capacity - Total) return false;
        Total += amount;
        return true;
    }
}
