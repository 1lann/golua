function init()
	local pid
	local cid = 0
	while true do
		if cid == 0 then
			if t0 then
				cid, pid = 2, 0
			else
				cid, pid = 1, 0
			end
		elseif cid == 1 then
			init()
			cid, pid = 2, 1
		elseif cid == 2 then
			return
		end
	end
end

function pi(n)
	local pid, t2, t5, t7, t9, t6, t8, t0, t1, t3, t4
	local cid = 0
	while true do
		if cid == 0 then
			cid, pid = 3, 0
		elseif cid == 1 then
			t0 = -1 ^ t8
			t1 = 4 * t0
			t2 = 2 * t8
			t3 = t2 + 1
			t4 = t1 / t3
			t5 = t7 + t4
			t6 = t8 + 1
			cid, pid = 3, 1
		elseif cid == 2 then
			return t7
		elseif cid == 3 then
			t7 = ({[0] = 0, [1] = t5})[pid]
			t8 = ({[0] = 0, [1] = t6})[pid]
			t9 = t8 < n
			if t9 then
				cid, pid = 1, 3
			else
				cid, pid = 2, 3
			end
		end
	end
end

function GCD(a,b)
	local pid, t0, t1, t2, t3
	local cid = 0
	while true do
		if cid == 0 then
			cid, pid = 3, 0
		elseif cid == 1 then
			t0 = t1 % t2
			cid, pid = 3, 1
		elseif cid == 2 then
			return t1
		elseif cid == 3 then
			t1 = ({[0] = a, [1] = t2})[pid]
			t2 = ({[0] = b, [1] = t0})[pid]
			t3 = t2 ~= 0
			if t3 then
				cid, pid = 1, 3
			else
				cid, pid = 2, 3
			end
		end
	end
end

function Phi(n)
	local pid, t4, t5, t6, t7, t0, t1, t2, t3
	local cid = 0
	while true do
		if cid == 0 then
			cid, pid = 3, 0
		elseif cid == 1 then
			t0 = GCD(t3,n)
			t1 = t0 == 1
			if t1 then
				cid, pid = 4, 1
			else
				cid, pid = 5, 1
			end
		elseif cid == 2 then
			return t2
		elseif cid == 3 then
			t2 = ({[0] = 1, [5] = t6})[pid]
			t3 = ({[0] = 2, [5] = t7})[pid]
			t4 = t3 < n
			if t4 then
				cid, pid = 1, 3
			else
				cid, pid = 2, 3
			end
		elseif cid == 4 then
			t5 = t2 + 1
			cid, pid = 5, 4
		elseif cid == 5 then
			t6 = ({[1] = t2, [4] = t5})[pid]
			t7 = t3 + 1
			cid, pid = 3, 5
		end
	end
end

function main()
	local pid, t1, t3
	local cid = 0
	while true do
		if cid == 0 then
			print("Hello, World!")
			t1 = Phi(1337)
			print(t1)
			t3 = pi(10000)
			print(t3)
			return
		end
	end
end

-- init()
main()

