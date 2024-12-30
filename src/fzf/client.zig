const std = @import("std");

const FzfSession = struct {
    fzf_cmd: std.process.Child,
    allocator: std.mem.Allocator,

    pub fn exec(self: *FzfSession, input: []const []const u8) ![]const u8 {
        self.fzf_cmd.stdin_behavior = .Pipe;
        self.fzf_cmd.stdout_behavior = .Pipe;
        _ = try self.fzf_cmd.spawn();

        _ = try writeInput(input, self.fzf_cmd);
        return try readOutput(self.fzf_cmd);
    }
};

/// Spawns fzf passes the input and returns the output once fzf is done running
pub fn exec(allocator: std.mem.Allocator, input: []const []const u8) ![]const u8 {
    var fzf_session = try init(allocator);
    return try fzf_session.exec(input);
}

fn init(allocator: std.mem.Allocator) !FzfSession {
    const fzf_cmd = std.process.Child.init(&[_][]const u8{"fzf"}, allocator);
    return FzfSession{ .fzf_cmd = fzf_cmd, .allocator = allocator };
}

fn writeInput(input: []const []const u8, cmd: std.process.Child) !void {
    const stdin = cmd.stdin orelse return error.FzfStdinNotFound;
    defer stdin.close();
    var writer = stdin.writer();
    for (input) |item| {
        writer.print("{s}\n", .{item}) catch |err| return err;
    }
}

fn readOutput(cmd: std.process.Child) ![]const u8 {
    const stdout = cmd.stdout orelse return error.FzfStdoutNotFound;
    defer stdout.close();
    // caution! read it only the end byte and no more since it can easily override the buffer
    const output_buffer = try stdout.readToEndAlloc(cmd.allocator, 512);
    return std.mem.trim(u8, output_buffer, &[_]u8{ 0, '\n' });
}
