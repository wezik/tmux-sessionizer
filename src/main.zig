const std = @import("std");
const persistence = @import("persistence/main.zig");

pub fn main() !void {
    var args = std.process.args();
    // skip program name
    _ = args.skip();

    const origin = args.next() orelse return;
    std.debug.print("origin: {s}\n", .{origin});
    const cmd = args.next() orelse return;
    std.debug.print("cmd: {s}\n", .{cmd});
    const Case = enum { a, add, c, create, h, help, v, version, l, list, r, remove, d, delete };
    const case = std.meta.stringToEnum(Case, cmd) orelse return;

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    switch (case) {
        .a, .add, .c, .create => {
            try create(origin, allocator);
        },
        .h, .help => help(),
        .v, .version => version(),
        .l, .list => {
            const config = try persistence.fetchConfig(allocator);
            for (config.entries) |entry| {
                std.debug.print("Config:\n  session_name: {s}\n  session_path: {s}\n  windows: {s}\n", .{
                    entry.session_name,
                    entry.session_path,
                    entry.windows,
                });
            }
        },
        .r, .remove, .d, .delete => delete(),
    }
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

pub fn delete() void {
    std.debug.print("delete\n", .{});
}

test "simple test" {}
