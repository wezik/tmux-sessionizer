const std = @import("std");
const tmux = @import("tmux/main.zig");
const fzf = @import("fzf-client.zig");

const ValidCommands = enum {
    s,
    select,
    a,
    add,
    c,
    create,
    h,
    help,
    v,
    version,
    l,
    list,
    r,
    remove,
    d,
    delete,
    e,
    edit,
};

/// Signals that can be sent to shell_script to perform actions
const Signal = enum {
    /// tmux attach signal
    tmux_attach,
    /// edit signal
    edit,

    pub fn send(self: Signal, string: []const u8) void {
        std.debug.print("signal:{s}:{s}\n", .{ @tagName(self), string });
    }
};

fn popArg(args: *std.process.ArgIterator) ValidCommands {
    // if no arg is passed, select is the case
    const cmd = args.next() orelse return ValidCommands.select;
    return std.meta.stringToEnum(ValidCommands, cmd) orelse return ValidCommands.select;
}

pub fn main() !void {
    var args = std.process.args();

    // skip program name
    _ = args.skip();
    const origin = args.next() orelse return;

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    switch (popArg(&args)) {
        .s, .select => try select(allocator),
        .a, .add, .c, .create => try create(origin, allocator),
        .h, .help => help(),
        .v, .version => version(),
        .l, .list => {
            const sessions = try tmux.getSessions(allocator);
            for (sessions) |entry| {
                std.debug.print("Config:\n  session_name: {s}\n  session_path: {s}\n  panes: {any}\n", .{
                    entry.name,
                    entry.path,
                    entry.windows,
                });
            }
        },
        .e, .edit => Signal.edit.send(try tmux.getConfigPath(allocator)),
        .r, .remove, .d, .delete => try delete(origin, allocator),
    }
}

pub fn select(allocator: std.mem.Allocator) !void {
    // prepare fzf input
    const sessions = try tmux.getSessions(allocator);
    var input = std.ArrayList([]const u8).init(allocator);
    var entryToSession = std.StringHashMap(tmux.TmuxSession).init(allocator);
    for (sessions) |session| {
        const entry = try tmux.sessionToKey(allocator, session);
        _ = try entryToSession.put(entry, session);
        _ = try input.append(entry);
    }

    // run fzf
    const output = try fzf.exec(allocator, try input.toOwnedSlice());
    const resultSession = entryToSession.get(output) orelse return;

    _ = try tmux.prepareSession(resultSession);
    Signal.tmux_attach.send(resultSession.name);
    return;
}

pub fn create(origin: []const u8, allocator: std.mem.Allocator) !void {
    const session = tmux.TmuxSession{
        .name = origin,
        .path = origin,
    };
    _ = try tmux.appendSession(allocator, session);
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
    _ = try tmux.deleteSession(origin, allocator);
}

test "simple test" {}
