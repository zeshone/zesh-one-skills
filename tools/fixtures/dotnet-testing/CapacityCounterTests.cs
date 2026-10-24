using Xunit;

namespace DotnetTestingTrial;

public sealed class CapacityCounterTests
{
    [Theory]
    [InlineData(1, 9)]
    [InlineData(2, 10)]
    public void TryAdd_AllowedAmount_ReturnsTrueAndUpdatesTotal(int amount, int expected)
    {
        var counter = new CapacityCounter();

        bool accepted = counter.TryAdd(amount);

        Assert.True(accepted);
        Assert.Equal(expected, counter.Total);
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-1)]
    [InlineData(3)]
    [InlineData(int.MaxValue)]
    public void TryAdd_RejectedAmount_ReturnsFalseAndPreservesTotal(int amount)
    {
        var counter = new CapacityCounter();
        int before = counter.Total;

        bool accepted = counter.TryAdd(amount);

        Assert.False(accepted);
        Assert.Equal(before, counter.Total);
    }
}
