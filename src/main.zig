const std = @import("std");
const persistence = @import("persistence/main.zig");

pub fn main() !void {
    var args = std.process.args();
    // skip program name
    _ = args.skip();

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    const origin = args.next() orelse return;
    std.debug.print("origin: {s}\n", .{origin});
    const cmd = blk: {
        const arg = args.next() orelse "";
        if (arg.len > 0) {
            break :blk arg;
        } else {
            try select(allocator);
            return;
        }
    };

    std.debug.print("cmd: {s}\n", .{cmd});
    const Case = enum { a, add, c, create, h, help, v, version, l, list, r, remove, d, delete };
    const case = std.meta.stringToEnum(Case, cmd) orelse return;

    switch (case) {
        .a, .add, .c, .create => try create(origin, allocator),
        .h, .help => help(),
        .v, .version => version(),
        .l, .list => {
            const config = try persistence.fetchConfig(allocator);
            for (config.entries) |entry| {
                std.debug.print("Config:\n  session_name: {s}\n  session_path: {s}\n  panes: {any}\n", .{
                    entry.session_name,
                    entry.session_path,
                    entry.panes,
                });
            }
        },
        .r, .remove, .d, .delete => try delete(origin, allocator),
    }
}

pub fn select(allocator: std.mem.Allocator) !void {
    // run fzf
    var fzf_cmd = std.process.Child.init(&[_][]const u8{"fzf"}, allocator);
    fzf_cmd.stdin_behavior = .Pipe;
    fzf_cmd.stdout_behavior = .Pipe;
    _ = try fzf_cmd.spawn();

    // write input to fzf stdin
    const stdin = fzf_cmd.stdin.?;
    const config = try persistence.fetchConfig(allocator);
    {
        var writer = stdin.writer();
        for (config.entries) |entry| {
            try writer.print("{s}\n", .{entry.session_path});
        }

        // close stdin after writing all input
        stdin.close();
    }

    // read fzf stdout
    const stdout = fzf_cmd.stdout.?;
    var output_buffer: [1024]u8 = undefined;
    const bytes_read = try stdout.readAll(&output_buffer);
    stdout.close();

    // parse fzf output
    const output = std.mem.trim(u8, output_buffer[0..bytes_read], &[_]u8{ 0, '\n' });
    std.debug.print("output: {s}\n", .{output});

    // TODO: setup or attach (meaning return the session name) to the selected session
    return;
}

pub fn create(origin: []const u8, allocator: std.mem.Allocator) !void {
    // Create directory if it doesn't exist
    const config = try persistence.fetchConfig(allocator);
    var new_config = std.ArrayList(persistence.ConfigEntry).init(allocator);
    _ = try new_config.appendSlice(config.entries);
    const new_entry = persistence.newEntry(origin);
    _ = try new_config.append(new_entry);
    try persistence.saveConfig(allocator, persistence.Config{ .entries = try new_config.toOwnedSlice() });
}

pub fn help() void {
    std.debug.print("help\n", .{});
}

pub fn version() void {
    std.debug.print("version\n", .{});
}

pub fn list() void {
    std.debug.print("list\n", .{});
}

pub fn delete(origin: []const u8, allocator: std.mem.Allocator) !void {
    const config = try persistence.fetchConfig(allocator);
    var new_config = std.ArrayList(persistence.ConfigEntry).init(allocator);
    for (config.entries) |entry| {
        if (!std.mem.eql(u8, entry.session_path, origin)) {
            _ = try new_config.append(entry);
        }
    }
    try persistence.saveConfig(allocator, persistence.Config{ .entries = try new_config.toOwnedSlice() });
}

test "simple test" {}
